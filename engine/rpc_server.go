package engine

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/funkygao/log4go"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// thrift.TServer implementation
type TFunServer struct {
	stopped bool

	engine                 *Engine
	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory

	pool *rpcThreadPool

	sessionN            int64 // concurrent sessions
	mu                  sync.Mutex
	clientConcurrencies map[string]int
}

func NewTFunServer(engine *Engine,
	processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TFunServer {
	this := &TFunServer{
		engine:                 engine,
		processorFactory:       thrift.NewTProcessorFactory(processor),
		serverTransport:        serverTransport,
		inputTransportFactory:  transportFactory,
		outputTransportFactory: transportFactory,
		inputProtocolFactory:   protocolFactory,
		outputProtocolFactory:  protocolFactory,
		clientConcurrencies:    make(map[string]int),
	}
	this.pool = newRpcThreadPool(this.engine.conf.rpc.pm, this.handleSession)
	engine.rpcThreadPool = this.pool

	// start the thread pool
	this.pool.Start()

	// any web frontend got stuck?
	go this.monitorClients()

	return this
}

func (this *TFunServer) Serve() error {
	this.stopped = false
	err := this.serverTransport.Listen()
	if err != nil {
		return err
	}

	for !this.stopped {
		client, err := this.serverTransport.Accept()
		if client != nil {
			this.pool.Dispatch(client)
		}

		if err != nil {
			log.Error("Accept: %v", err)
		}
	}

	return errors.New("rpc server stopped")
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
		if !this.engine.conf.rpc.tcpNoDelay {
			// golang is tcp no delay by default
			tcpClient.SetNoDelay(false)
		}

		log.Trace("session[%s] open", tcpClient.RemoteAddr())

		// store client concurrent connections count TODO
		this.mu.Lock()
		p := strings.SplitN(tcpClient.RemoteAddr().String(), ":", 2)
		if len(p) == 2 && p[0] != "" {
			this.clientConcurrencies[p[0]] += 1
			defer func() {
				this.mu.Lock()
				this.clientConcurrencies[p[0]] -= 1
				this.mu.Unlock()
			}()
		}
		this.mu.Unlock()
	} else {
		log.Error("non tcp conn found, should never happen")
		return
	}

	t1 := time.Now()
	remoteAddr := transport.(*thrift.TSocket).Conn().RemoteAddr().String()
	if err := this.processRequests(transport); err != nil {
		this.engine.stats.TotalFailedSessions.Inc(1)
		log.Error("session[%s]: %s", remoteAddr, err.Error())
	}

	elapsed := time.Since(t1)
	this.engine.stats.SessionLatencies.Update(elapsed.Nanoseconds() / 1e6)
	log.Trace("session[%s] close in %s", remoteAddr, elapsed)

	if elapsed > this.engine.conf.rpc.sessionSlowThreshold {
		this.engine.stats.TotalSlowSessions.Inc(1)
		log.Warn("session[%s] SLOW %s", remoteAddr, elapsed)
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

// if concurrent conns from same client is too high, it means
// web frontend(php-fpm) got stuck, keep forking children
// TODO
func (this *TFunServer) monitorClients() {
	for {
		for clientIp, concurrentConns := range this.clientConcurrencies {
			if concurrentConns > 100 { // TODO
				log.Warn("Client[%s] may got stuck: %d", clientIp, concurrentConns)
			}
		}

		time.Sleep(time.Second * 10)
	}
}

func (this *TFunServer) Stop() error {
	this.stopped = true
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
