package engine

import (
	log "code.google.com/p/log4go"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"os"
	"time"
)

func (this *Engine) ServeForever() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()

	this.launchHttpServ()
	defer this.stopHttpServ()

	done := make(chan int)
	go this.launchRpcServe(done)
	<-done
}

func (this *Engine) launchRpcServe(done chan<- int) {
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
		fmt.Fprintf(os.Stderr, "Invalid protocol: %s\n", this.conf.rpc.protocol)
		os.Exit(1)
	}

	transportFactory := thrift.NewTTransportFactory()
	if this.conf.rpc.framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	serverTransport, err := thrift.NewTServerSocket(this.conf.rpc.listenAddr)
	if err != nil {
		panic(err)
	}

	rpcServer := thrift.NewTSimpleServer4(this.rpcProcessor,
		serverTransport, transportFactory, protocolFactory)
	log.Info("RPC server ready at %s", this.conf.rpc.listenAddr)

	for {
		err = rpcServer.Serve()
		if err != nil {
			log.Error(err)
			break
		}
	}

	done <- 1
}
