package engine

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"strings"
)

// thrift internal layer
//
// Server
// Processor (compiler genereated)
// Protocol (JSON/compact/...)
// Transport (TCP/HTTP/...)
func (this *Engine) launchRpcServe() (done chan interface{}) {
	var (
		protocolFactory  thrift.TProtocolFactory
		serverTransport  thrift.TServerTransport
		transportFactory thrift.TTransportFactory
		err              error
		serverNetwork    string
	)

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

	switch {
	case this.conf.rpc.framed:
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	default:
		transportFactory = thrift.NewTBufferedTransportFactory(2 << 10) // TODO
	}

	switch {
	case strings.Contains(this.conf.rpc.listenAddr, "/"):
		serverNetwork = "unix"
		serverTransport, err = NewTUnixSocketTimeout(
			this.conf.rpc.listenAddr, this.conf.rpc.sessionTimeout)

	default:
		serverNetwork = "tcp"
		serverTransport, err = thrift.NewTServerSocketTimeout(
			this.conf.rpc.listenAddr, this.conf.rpc.sessionTimeout)
	}
	if err != nil {
		panic(err)
	}

	// when config loaded, create the servants
	svr := servant.NewFunServant(config.Servants)
	this.rpcProcessor = rpc.NewFunServantProcessor(svr)
	svr.Start()

	this.rpcServer = NewTFunServer(this, this.rpcProcessor,
		serverTransport, transportFactory, protocolFactory)
	log.Info("RPC server ready at %s:%s", serverNetwork, this.conf.rpc.listenAddr)

	done = make(chan interface{})
	go func() {
		for {
			err = this.rpcServer.Serve()
			if err != nil {
				log.Error("rpcServer: %+v", err)
				break
			}
		}

		done <- 1

	}()

	return done
}
