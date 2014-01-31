package engine

import (
	log "code.google.com/p/log4go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fxi/servant"
	"github.com/funkygao/fxi/servant/gen-go/fun/rpc"
)

func (this *Engine) ServeForever() {
	handler := servant.NewFunServant()
	processor := rpc.NewFunServantProcessor(handler)
	listenSocket, err := thrift.NewTServerSocket(this.conf.listenAddr)
	if err != nil {
		panic(err)
	}

	server := thrift.NewTSimpleServer2(processor, listenSocket)
	log.Info("Server ready at %s", this.conf.listenAddr)

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
