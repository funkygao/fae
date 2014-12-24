package engine

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/etclib"
	log "github.com/funkygao/log4go"
	"net"
	"sync/atomic"
	"time"
)

// thrift.TServer implementation
type TFunServer struct {
	quit chan bool

	engine                 *Engine
	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory

	pool *rpcThreadPool

	sessionN int64 // concurrent sessions
}

func NewTFunServer(engine *Engine,
	processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TFunServer {
	this := &TFunServer{
		quit:                   make(chan bool),
		engine:                 engine,
		processorFactory:       thrift.NewTProcessorFactory(processor),
		serverTransport:        serverTransport,
		inputTransportFactory:  transportFactory,
		outputTransportFactory: transportFactory,
		inputProtocolFactory:   protocolFactory,
		outputProtocolFactory:  protocolFactory,
	}
	this.pool = newRpcThreadPool(this.engine.conf.rpc.maxOutstandingSessions,
		this.handleSession)
	engine.rpcThreadPool = this.pool

	return this
}

func (this *TFunServer) Serve() error {
	const stoppedError = "RPC server stopped"

	err := this.serverTransport.Listen()
	if err != nil {
		return err
	}

	// register to etcd
	// once registered, other peers will connect to me
	// so, must be after Listen ready
	if this.engine.conf.EtcdSelfAddr != "" {
		etclib.BootService(this.engine.conf.EtcdSelfAddr, etclib.SERVICE_FAE)

		log.Info("etcd self[%s] registered", this.engine.conf.EtcdSelfAddr)
	}

	for {
		select {
		case <-this.quit:
			// FIXME new conn will timeout, instead of conn close
			log.Info("RPC server quit...")
			return errors.New(stoppedError)

		default:
		}

		client, err := this.serverTransport.Accept()
		if err != nil {
			log.Error("Accept: %v", err)
		} else {
			this.pool.Dispatch(client)
		}
	}

	return errors.New(stoppedError)
}

func (this *TFunServer) handleSession(client interface{}) {
	defer atomic.AddInt64(&this.sessionN, -1)

	transport, ok := client.(thrift.TTransport)
	if !ok {
		log.Error("Invalid client: %#v", client)
		return
	}

	this.engine.stats.SessionPerSecond.Mark(1)
	atomic.AddInt64(&this.sessionN, 1)

	if tcpClient, ok := transport.(*thrift.TSocket).Conn().(*net.TCPConn); ok {
		log.Trace("session[%s] open", tcpClient.RemoteAddr())
	} else {
		log.Error("non tcp conn found, should NEVER happen")
		return
	}

	t1 := time.Now()
	remoteAddr := transport.(*thrift.TSocket).Conn().RemoteAddr().String()
	if err := this.processRequests(transport); err != nil {
		this.engine.stats.TotalFailedSessions.Inc(1)
	}

	elapsed := time.Since(t1)
	this.engine.stats.SessionLatencies.Update(elapsed.Nanoseconds() / 1e6)
	log.Trace("session[%s] close in %s", remoteAddr, elapsed)

	if this.engine.conf.rpc.sessionSlowThreshold.Seconds() > 0 &&
		elapsed > this.engine.conf.rpc.sessionSlowThreshold {
		this.engine.stats.TotalSlowSessions.Inc(1)
	}
}

func (this *TFunServer) processRequests(client thrift.TTransport) error {
	processor := this.processorFactory.GetProcessor(client)
	inputTransport := this.inputTransportFactory.GetTransport(client)
	outputTransport := this.outputTransportFactory.GetTransport(client)
	inputProtocol := this.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := this.outputProtocolFactory.GetProtocol(outputTransport)
	defer func() {
		if inputTransport != nil {
			inputTransport.Close()
		}
		if outputTransport != nil {
			outputTransport.Close()
		}
	}()

	var (
		t1        time.Time
		elapsed   time.Duration
		tcpClient = client.(*thrift.TSocket).Conn().(*net.TCPConn)
		callsN    int64
	)

	for {
		t1 = time.Now()
		if this.engine.conf.rpc.ioTimeout > 0 { // read + write
			tcpClient.SetDeadline(time.Now().Add(this.engine.conf.rpc.ioTimeout))
		}

		ok, err := processor.Process(inputProtocol, outputProtocol)
		if err == nil {
			callsN++
		}

		elapsed = time.Since(t1)
		this.engine.stats.CallLatencies.Update(elapsed.Nanoseconds() / 1e6)
		this.engine.stats.CallPerSecond.Mark(1)

		// check transport error
		if err, ok := err.(thrift.TTransportException); ok &&
			err.TypeId() == thrift.END_OF_FILE {
			// remote client closed transport, this is normal end of session
			log.Trace("session[%s] %d calls EOF", tcpClient.RemoteAddr().String(),
				callsN)
			this.engine.stats.CallPerSession.Update(callsN)
			return nil
		} else if err != nil {
			// non-EOF transport err
			// e,g. connection reset by peer
			// e,g. broken pipe
			// e,g. read tcp i/o timeout
			this.engine.stats.TotalFailedCalls.Inc(1)
			this.engine.stats.CallPerSession.Update(callsN)

			log.Trace("session[%s] %d calls: %s",
				tcpClient.RemoteAddr().String(), callsN, err.Error())
			return err
		}

		// it is servant generated TApplicationException
		// err logging is handled inside servants
		if err != nil {
			this.engine.stats.TotalFailedCalls.Inc(1)

			log.Trace("session[%s] %d calls: %s",
				tcpClient.RemoteAddr().String(),
				callsN, err.Error())
		}

		if !ok || !inputProtocol.Transport().Peek() {
			break
		}
	}

	this.engine.stats.CallPerSession.Update(callsN)
	return nil
}

func (this *TFunServer) Stop() error {
	close(this.quit)
	this.serverTransport.Interrupt()
	return nil
}

func (this *TFunServer) ProcessorFactory() thrift.TProcessorFactory {
	return this.processorFactory
}

func (this *TFunServer) ServerTransport() thrift.TServerTransport {
	return this.serverTransport
}

func (this *TFunServer) InputTransportFactory() thrift.TTransportFactory {
	return this.inputTransportFactory
}

func (this *TFunServer) OutputTransportFactory() thrift.TTransportFactory {
	return this.outputTransportFactory
}

func (this *TFunServer) InputProtocolFactory() thrift.TProtocolFactory {
	return this.inputProtocolFactory
}

func (this *TFunServer) OutputProtocolFactory() thrift.TProtocolFactory {
	return this.outputProtocolFactory
}
