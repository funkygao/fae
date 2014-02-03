package thriftex

import (
	log "code.google.com/p/log4go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

// thrift.TServer implementation
type TSimpleServer struct {
	stopped bool

	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory
}

func NewTSimpleServer4(processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TSimpleServer {
	return NewTSimpleServerFactory4(thrift.NewTProcessorFactory(processor),
		serverTransport,
		transportFactory,
		protocolFactory,
	)
}

func NewTSimpleServerFactory4(processorFactory thrift.TProcessorFactory,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TSimpleServer {
	return NewTSimpleServerFactory6(processorFactory,
		serverTransport,
		transportFactory,
		transportFactory,
		protocolFactory,
		protocolFactory,
	)
}

func NewTSimpleServerFactory6(processorFactory thrift.TProcessorFactory,
	serverTransport thrift.TServerTransport,
	inputTransportFactory thrift.TTransportFactory,
	outputTransportFactory thrift.TTransportFactory,
	inputProtocolFactory thrift.TProtocolFactory,
	outputProtocolFactory thrift.TProtocolFactory) *TSimpleServer {
	return &TSimpleServer{processorFactory: processorFactory,
		serverTransport:        serverTransport,
		inputTransportFactory:  inputTransportFactory,
		outputTransportFactory: outputTransportFactory,
		inputProtocolFactory:   inputProtocolFactory,
		outputProtocolFactory:  outputProtocolFactory,
	}
}

func (p *TSimpleServer) ProcessorFactory() thrift.TProcessorFactory {
	return p.processorFactory
}

func (p *TSimpleServer) ServerTransport() thrift.TServerTransport {
	return p.serverTransport
}

func (p *TSimpleServer) InputTransportFactory() thrift.TTransportFactory {
	return p.inputTransportFactory
}

func (p *TSimpleServer) OutputTransportFactory() thrift.TTransportFactory {
	return p.outputTransportFactory
}

func (p *TSimpleServer) InputProtocolFactory() thrift.TProtocolFactory {
	return p.inputProtocolFactory
}

func (p *TSimpleServer) OutputProtocolFactory() thrift.TProtocolFactory {
	return p.outputProtocolFactory
}

func (p *TSimpleServer) Serve() error {
	p.stopped = false
	err := p.serverTransport.Listen()
	if err != nil {
		return err
	}
	for !p.stopped {
		client, err := p.serverTransport.Accept()
		if err != nil {
			log.Error("Accept err: ", err)
		}

		if client != nil {
			go func() {
				if err := p.processRequest(client); err != nil {
					log.Error("error processing request:", err)
				}
			}()
		}
	}
	return nil
}

func (p *TSimpleServer) Stop() error {
	p.stopped = true
	p.serverTransport.Interrupt()
	return nil
}

func (p *TSimpleServer) processRequest(client thrift.TTransport) error {
	processor := p.processorFactory.GetProcessor(client)
	inputTransport := p.inputTransportFactory.GetTransport(client)
	outputTransport := p.outputTransportFactory.GetTransport(client)
	inputProtocol := p.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := p.outputProtocolFactory.GetProtocol(outputTransport)
	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}
	for {
		ok, err := processor.Process(inputProtocol, outputProtocol)
		if err, ok := err.(thrift.TTransportException); ok && err.TypeId() == thrift.END_OF_FILE {
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
