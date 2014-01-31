package engine

import (
	log "code.google.com/p/log4go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fxi/servant"
	"github.com/funkygao/fxi/servant/gen-go/fun/rpc"
	"os"
	"time"
)

func (this *Engine) ServeForever() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()

	this.launchHttpServ()
	defer this.stopHttpServ()

	handler := servant.NewFunServant()
	processor := rpc.NewFunServantProcessor(handler)
	listenSocket, err := thrift.NewTServerSocket(this.conf.rpcListenAddr)
	if err != nil {
		panic(err)
	}

	rpcServer := thrift.NewTSimpleServer2(processor, listenSocket)
	log.Info("RPC server ready at %s", this.conf.rpcListenAddr)

	if err := rpcServer.Serve(); err != nil {
		panic(err)
	}
}
