package engine

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
	"net"
	"sync/atomic"
	"time"
)

// thrift.TServer implementation
type TFunServer struct {
	quit           chan bool
	activeSessionN int64

	engine                 *Engine
	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory

	pool *rpcThreadPool
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
	this.pool = newRpcThreadPool(
		config.Engine.Rpc.MaxOutstandingSessions,
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
	if config.Engine.EtcdSelfAddr != "" {
		etclib.BootService(config.Engine.EtcdSelfAddr, etclib.SERVICE_FAE)

		log.Info("etcd self[%s] registered", config.Engine.EtcdSelfAddr)
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
	defer atomic.AddInt64(&this.activeSessionN, -1)

	transport, ok := client.(thrift.TTransport)
	if !ok {
		// should never happen
		log.Error("Invalid client: %#v", client)
		return
	}

	this.engine.stats.SessionPerSecond.Mark(1)
	atomic.AddInt64(&this.activeSessionN, 1)

	if tcpClient, ok := transport.(*thrift.TSocket).Conn().(*net.TCPConn); ok {
		log.Debug("session[%s] open", tcpClient.RemoteAddr())
	} else {
		log.Error("non tcp conn found, should NEVER happen")
		return
	}

	t1 := time.Now()
	remoteAddr := transport.(*thrift.TSocket).Conn().RemoteAddr().String()
	var (
		calls int64
		err   error
	)
	if calls, err = this.processRequests(transport); err != nil {
		this.engine.stats.TotalFailedSessions.Inc(1)
	}

	elapsed := time.Since(t1)
	this.engine.stats.SessionLatencies.Update(elapsed.Nanoseconds() / 1e6)
	if err != nil {
		log.Error("session[%s] %d calls in %s: %v", remoteAddr, calls, elapsed, err)
	} else {
		log.Trace("session[%s] %d calls in %s", remoteAddr, calls, elapsed)
	}

	if config.Engine.Rpc.SessionSlowThreshold.Seconds() > 0 &&
		elapsed > config.Engine.Rpc.SessionSlowThreshold {
		this.engine.stats.TotalSlowSessions.Inc(1)
	}
}

// FIXME if error occurs, fae will actively close this session: halt on first exception
func (this *TFunServer) processRequests(client thrift.TTransport) (int64, error) {
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
		if config.Engine.Rpc.IoTimeout > 0 { // read + write
			tcpClient.SetDeadline(time.Now().Add(config.Engine.Rpc.IoTimeout))
		}

		ok, ex := processor.Process(inputProtocol, outputProtocol)
		if ex == nil {
			callsN++
		}

		elapsed = time.Since(t1)
		this.engine.stats.CallLatencies.Update(elapsed.Nanoseconds() / 1e6)
		this.engine.stats.CallPerSecond.Mark(1)

		// check transport error
		if err, isTransportEx := ex.(thrift.TTransportException); isTransportEx &&
			err.TypeId() == thrift.END_OF_FILE {
			// remote client closed transport, this is normal end of session
			this.engine.stats.CallPerSession.Update(callsN)
			return callsN, nil
		} else if err != nil {
			// non-EOF transport err
			// e,g. connection reset by peer
			// e,g. broken pipe
			// e,g. read tcp i/o timeout
			this.engine.stats.TotalFailedCalls.Inc(1)
			this.engine.stats.CallPerSession.Update(callsN)

			return callsN, err
		}

		// it is servant generated TApplicationException
		// e,g Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'WHERE entityId=?' at line 1
		if ex != nil {
			this.engine.stats.TotalFailedCalls.Inc(1)
			callsN++
			return callsN, ex // TODO stop the session?
		}

		// Peek: there is more data to be read or the remote side is still open
		if !ok || !inputProtocol.Transport().Peek() {
			break // TODO stop the session?
		}
	}

	this.engine.stats.CallPerSession.Update(callsN)
	return callsN, nil
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
