package engine

import (
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

type Engine struct {
	conf *Config

	StartedAt time.Time

	listener   net.Listener
	httpServer *http.Server
	httpRouter *mux.Router
	httpPaths  []string

	stats    *engineStats
	pid      int
	hostname string
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	this.stats = newEngineStats(this)
	return
}
