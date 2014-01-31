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
	stats     *engineStats

	listener   net.Listener
	httpServer *http.Server
	httpRouter *mux.Router
	httpPaths  []string

	pid      int
	hostname string
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	this.stats = newEngineStats(this)
	return
}
