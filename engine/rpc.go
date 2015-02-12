package engine

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"strings"
	"sync/atomic"
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

	switch config.Engine.Rpc.Protocol {
	case "binary":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()

	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()

	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()

	default:
		panic(fmt.Sprintf("Invalid protocol: %s", config.Engine.Rpc.Protocol))
	}

	// client-side Thrift protocol/transport stack must match
	// the server-side, otherwise you are very likely to get in trouble
	switch {
	case config.Engine.Rpc.Framed:
		// each payload is sent over the wire with a frame header containing its size
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	default:
		// there is no BufferedTransport in Java: only FramedTransport
		transportFactory = thrift.NewTBufferedTransportFactory(
			config.Engine.Rpc.BufferSize)
	}

	switch {
	case strings.Contains(config.Engine.Rpc.ListenAddr, "/"):
		serverNetwork = "unix"
		if config.Engine.Rpc.SessionTimeout > 0 {
			serverTransport, err = NewTUnixSocketTimeout(
				config.Engine.Rpc.ListenAddr, config.Engine.Rpc.SessionTimeout)
		} else {
			serverTransport, err = NewTUnixSocket(
				config.Engine.Rpc.ListenAddr)
		}

	default:
		serverNetwork = "tcp"
		if config.Engine.Rpc.SessionTimeout > 0 {
			serverTransport, err = thrift.NewTServerSocketTimeout(
				config.Engine.Rpc.ListenAddr, config.Engine.Rpc.SessionTimeout)
		} else {
			serverTransport, err = thrift.NewTServerSocket(
				config.Engine.Rpc.ListenAddr)
		}
	}
	if err != nil {
		panic(err)
	}

	// dial zk before startup servants
	// because proxy servant is dependent upon zk
	if config.Engine.EtcdSelfAddr != "" {
		if err := etclib.Dial(config.Engine.EtcdServers); err != nil {
			panic(err)
		} else {
			log.Debug("etcd connected: %+v", config.Engine.EtcdServers)
		}
	}

	// when config loaded, create the servants
	this.svt = servant.NewFunServant(config.Engine.Servants)
	this.rpcProcessor = rpc.NewFunServantProcessor(this.svt)
	this.svt.Start()

	this.rpcServer = NewTFunServer(this, this.rpcProcessor,
		serverTransport, transportFactory, protocolFactory)
	log.Info("RPC server ready at %s:%s", serverNetwork, config.Engine.Rpc.ListenAddr)

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

func (this *Engine) StopRpcServe() {
	rpcServer := this.rpcServer.(*TFunServer)
	rpcServer.Stop()

	close(this.stopChan)

	outstandingSessions := atomic.LoadInt64(&rpcServer.activeSessionN)
	log.Warn("RPC outstanding sessions: %d", outstandingSessions)

	this.svt.Flush()

	log.Info("RPC server stopped gracefully")
}
