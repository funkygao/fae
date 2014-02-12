package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/funkygao/log4go"
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
}

func NewTFunServer(engine *Engine,
	processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TFunServer {
	return &TFunServer{
		engine:                 engine,
		processorFactory:       thrift.NewTProcessorFactory(processor),
		serverTransport:        serverTransport,
		inputTransportFactory:  transportFactory,
		outputTransportFactory: transportFactory,
		inputProtocolFactory:   protocolFactory,
		outputProtocolFactory:  protocolFactory,
	}
}

func (this *TFunServer) Serve() error {
	this.stopped = false
	err := this.serverTransport.Listen()
	if err != nil {
		return err
	}

	for !this.stopped {
		client, err := this.serverTransport.Accept()
		if err != nil {
			log.Error("Accept err: ", err)
		}

		if client != nil {
			this.engine.stats.TotalSessions.Add(1)
			if this.engine.conf.rpc.debugSession {
				log.Debug("accepted session peer %s",
					client.(*thrift.TSocket).Conn().RemoteAddr().String())
			}

			go this.processSession(client)
		}
	}

	return nil
}

func (this *TFunServer) processSession(client thrift.TTransport) {
	t1 := time.Now()
	remoteAddr := client.(*thrift.TSocket).Conn().RemoteAddr().String()
	if err := this.processRequest(client); err != nil {
		this.engine.stats.TotalFailedSessions.Add(1)
		log.Error("session peer[%s] failed: %s", remoteAddr, err)
	}

	if this.engine.conf.rpc.debugSession {
		log.Debug("session peer[%s] closed", remoteAddr)
	}

	elapsed := time.Since(t1)
	if elapsed.Seconds() > this.engine.conf.rpc.sessionSlowThreshold {
		// slow session
		this.engine.stats.TotalSlowSessions.Add(1)
		log.Warn("SLOW=%s session peer: %s", elapsed, remoteAddr)
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

		this.engine.stats.TotalCalls.Add(1)
		elapsed = time.Since(t1)
		if elapsed.Seconds() > this.engine.conf.rpc.callSlowThreshold {
			// slow call
			this.engine.stats.TotalSlowCalls.Add(1)
			log.Warn("SLOW=%s call peer: %s", elapsed, remoteAddr)
		}

		// check transport error
		if err, ok := err.(thrift.TTransportException); ok &&
			err.TypeId() == thrift.END_OF_FILE {
			// remote client closed transport
			return nil
		} else if err != nil {
			// non-EOF transport err
			log.Error("ERROR transport, peer(%s) %s", remoteAddr, err)
			this.engine.stats.TotalFailedCalls.Add(1)
			return err
		}

		// it is servant generated TApplicationException
		if err != nil {
			this.engine.stats.TotalFailedCalls.Add(1)
			log.Error("ERROR servant call, peer(%s) %s", remoteAddr, err)
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
