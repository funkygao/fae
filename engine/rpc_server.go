package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/funkygao/log4go"
	"net"
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
	}
	this.pool = newRpcThreadPool(this.engine.conf.rpc.pm, this.handleClient)
	return this
}

func (this *TFunServer) Serve() error {
	this.stopped = false
	err := this.serverTransport.Listen()
	if err != nil {
		return err
	}

	this.pool.start()

	for !this.stopped {
		client, err := this.serverTransport.Accept()
		if client != nil {
			this.pool.dispatch(client)
		}

		if err != nil {
			log.Error("Accept: %v", err)
		}
	}

	return nil
}

func (this *TFunServer) handleClient(req interface{}) {
	defer func() {
		this.engine.stats.CurrentSessions.Dec(1)
	}()

	client := req.(thrift.TTransport)

	this.engine.stats.SessionPerSecond.Mark(1)
	this.engine.stats.CurrentSessions.Inc(1)

	if tcp, ok := client.(*thrift.TSocket).Conn().(*net.TCPConn); ok {
		tcp.SetNoDelay(this.engine.conf.rpc.tcpNoDelay)

		if this.engine.conf.rpc.debugSession {
			log.Debug("accepted session peer{%s}", tcp.RemoteAddr())
		}
	}

	this.processSession(client)
}

func (this *TFunServer) processSession(client thrift.TTransport) {
	t1 := time.Now()
	remoteAddr := client.(*thrift.TSocket).Conn().RemoteAddr().String()
	if err := this.processRequest(client); err != nil {
		this.engine.stats.TotalFailedSessions.Inc(1)
		log.Error("session peer{%s}: %s", remoteAddr, err)
	}

	elapsed := time.Since(t1)
	this.engine.stats.SessionLatencies.Update(elapsed.Nanoseconds() / 1e6)
	if this.engine.conf.rpc.debugSession {
		log.Debug("session peer{%s} closed after %s", remoteAddr, elapsed)
	} else if elapsed.Seconds() > this.engine.conf.rpc.sessionSlowThreshold {
		// slow session
		this.engine.stats.TotalSlowSessions.Inc(1)
		log.Warn("SLOW=%s session peer{%s}", elapsed, remoteAddr)
	}

}

func (this *TFunServer) processRequest(client thrift.TTransport) error {
	processor := this.processorFactory.GetProcessor(client)
	inputTransport := this.inputTransportFactory.GetTransport(client)
	outputTransport := this.outputTransportFactory.GetTransport(client)
	inputProtocol := this.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := this.outputProtocolFactory.GetProtocol(outputTransport)
	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}

	var (
		t1         time.Time
		elapsed    time.Duration
		remoteAddr = client.(*thrift.TSocket).Conn().RemoteAddr().String()
	)
	for {
		t1 = time.Now()
		ok, err := processor.Process(inputProtocol, outputProtocol)

		elapsed = time.Since(t1)
		this.engine.stats.CallPerSecond.Mark(1)
		this.engine.stats.CallLatencies.Update(elapsed.Nanoseconds() / 1e6)
		if elapsed.Seconds() > this.engine.conf.rpc.callSlowThreshold {
			// slow call
			this.engine.stats.TotalSlowCalls.Inc(1)
			log.Warn("SLOW call=%.3fs, peer{%s}", elapsed.Seconds(), remoteAddr)
		}

		// check transport error
		if err, ok := err.(thrift.TTransportException); ok &&
			err.TypeId() == thrift.END_OF_FILE {
			// remote client closed transport
			return nil
		} else if err != nil {
			// non-EOF transport err
			// e,g. connection reset by peer
			// e,g. broken pipe
			this.engine.stats.TotalFailedCalls.Inc(1)
			return err
		}

		// it is servant generated TApplicationException
		if err != nil {
			this.engine.stats.TotalFailedCalls.Inc(1)
			log.Error("servant call peer{%s}: %s", remoteAddr, err)
		}

		if !ok || !inputProtocol.Transport().Peek() {
			break
		}
	}

	return nil
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
