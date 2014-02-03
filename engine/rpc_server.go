package engine

import (
	log "code.google.com/p/log4go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

// thrift.TServer implementation
type TFunServer struct {
	stopped bool

	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory
}

func NewTFunServer(processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TFunServer {
	return &TFunServer{
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
			//log.Debug("%+v", client.(thrift.TSocket).Conn().RemoteAddr())

			go func() {
				if err := this.processRequest(client); err != nil {
					log.Error("error processing request:", err)
				}
			}()
		}
	}

	return nil
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

	for {
		ok, err := processor.Process(inputProtocol, outputProtocol)
		if err, ok := err.(thrift.TTransportException); ok &&
			err.TypeId() == thrift.END_OF_FILE {
			return nil
		} else if err != nil {
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
