package engine

import (
	log "code.google.com/p/log4go"
	"git.apache.org/thrift.git/lib/go/thrift"
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

		log.Debug("new client %v", client.(*thrift.TSocket).Conn().RemoteAddr())

		if client != nil {
			this.engine.stats.TotalSessions.Add(1)
			go this.processSession(client)
		}
	}

	return nil
}

func (this *TFunServer) processSession(client thrift.TTransport) {
	t1 := time.Now()
	if err := this.processRequest(client); err != nil {
		log.Error("error processing request: ", err)
	}

	elapsed := time.Since(t1)
	if elapsed.Seconds() > this.engine.conf.rpc.clientSlowThreshold {
		// slow query
		log.Warn("client closed after %s", elapsed)
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
		t1      time.Time
		elapsed time.Duration
	)
	for {
		t1 = time.Now()
		ok, err := processor.Process(inputProtocol, outputProtocol)
		this.engine.stats.TotalCalls.Add(1)

		elapsed = time.Since(t1)
		if elapsed.Seconds() > this.engine.conf.rpc.callSlowThreshold {
			// slow query
			log.Warn("processed in %s", elapsed)
		}

		if err, ok := err.(thrift.TTransportException); ok &&
			err.TypeId() == thrift.END_OF_FILE {
			return nil
		} else if err != nil {
			this.engine.stats.TotalFailedCalls.Add(1)
			return err
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
