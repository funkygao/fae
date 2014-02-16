package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

type Engine struct {
	conf *engineConfig

	configFile string
	StartedAt  time.Time

	rpcProcessor thrift.TProcessor
	rpcServer    thrift.TServer

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
