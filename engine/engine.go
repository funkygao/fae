package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

type Engine struct {
	conf *engineConfig

	configFile string
	StartedAt  time.Time

	httpListener net.Listener
	httpServer   *http.Server
	httpRouter   *mux.Router
	httpPaths    []string

	rpcProcessor *rpc.FunServantProcessor
	rpcServer    *thrift.TSimpleServer

	stats    *engineStats
	pid      int
	hostname string
}

func NewEngine(fn string) (this *Engine) {
	this = new(Engine)
	this.configFile = fn
	this.stats = newEngineStats(this)

	return
}
