package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
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

	rpcProcessor thrift.TProcessor
	rpcServer    thrift.TServer

	peer *Peer

	stats    *engineStats
	pid      int
	hostname string
}

func NewEngine(fn string) (this *Engine) {
	this = new(Engine)
	this.configFile = fn
	this.stats = newEngineStats()

	return
}
