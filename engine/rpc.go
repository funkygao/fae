package engine

import (
	log "code.google.com/p/log4go"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"strings"
)

func (this *Engine) launchRpcServe() (done chan interface{}) {
	var protocolFactory thrift.TProtocolFactory
	switch this.conf.rpc.protocol {
	case "binary":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()

	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()

	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()

	default:
		panic(fmt.Sprintf("Invalid protocol: %s", this.conf.rpc.protocol))
	}

	transportFactory := thrift.NewTTransportFactory()
	if this.conf.rpc.framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	var (
		serverTransport thrift.TServerTransport
		err             error
	)
	switch {
	case strings.Contains(this.conf.rpc.listenAddr, "/"):
		serverTransport, err = NewTUnixSocketTimeout(
			this.conf.rpc.listenAddr, this.conf.rpc.clientTimeout)

	default:
		serverTransport, err = thrift.NewTServerSocketTimeout(
			this.conf.rpc.listenAddr, this.conf.rpc.clientTimeout)
	}
	if err != nil {
		panic(err)
	}

	this.rpcServer = NewTFunServer(this, this.rpcProcessor,
		serverTransport, transportFactory, protocolFactory)
	log.Info("RPC server ready at %s", this.conf.rpc.listenAddr)

	done = make(chan interface{})
	go func() {
		for {
			err = this.rpcServer.Serve()
			if err != nil {
				log.Error(err)
				break
			}
		}

		done <- 1

	}()

	return done
}
