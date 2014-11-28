package engine

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"strings"
)

// thrift internal layer
//
// Server
// Processor (compiler genereated)
// Protocol (JSON/compact/...), what is transmitted
// Transport (TCP/HTTP/...), how is transmitted
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

	// client-side Thrift protocol/transport stack must match
	// the server-side, otherwise you are very likely to get in trouble
	switch {
	case this.conf.rpc.framed:
		// each payload is sent over the wire with a frame header containing its size
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	default:
		// there is no BufferedTransport in Java: only FramedTransport
		transportFactory = thrift.NewTBufferedTransportFactory(this.conf.rpc.bufferSize)
	}

	switch {
	case strings.Contains(this.conf.rpc.listenAddr, "/"):
		serverNetwork = "unix"
		if this.conf.rpc.sessionTimeout.Seconds() > 0 {
			serverTransport, err = NewTUnixSocketTimeout(
				this.conf.rpc.listenAddr, this.conf.rpc.sessionTimeout)
		} else {
			serverTransport, err = NewTUnixSocket(
				this.conf.rpc.listenAddr)
		}

	default:
		serverNetwork = "tcp"
		if this.conf.rpc.sessionTimeout.Seconds() > 0 {
			serverTransport, err = thrift.NewTServerSocketTimeout(
				this.conf.rpc.listenAddr, this.conf.rpc.sessionTimeout)
		} else {
			serverTransport, err = thrift.NewTServerSocket(
				this.conf.rpc.listenAddr)
		}
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
