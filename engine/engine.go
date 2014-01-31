package engine

import (
	"github.com/funkygao/fxi/servant"
	"github.com/funkygao/fxi/servant/gen-go/fun/rpc"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

type Engine struct {
	conf *Config

	StartedAt time.Time

	httpListener net.Listener
	httpServer   *http.Server
	httpRouter   *mux.Router
	httpPaths    []string

	rpcProcessor *rpc.FunServantProcessor

	stats    *engineStats
	pid      int
	hostname string
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	this.stats = newEngineStats(this)
	this.rpcProcessor = rpc.NewFunServantProcessor(servant.NewFunServant())
	return
}
