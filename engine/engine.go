package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

type Engine struct {
	conf *engineConfig

	StartedAt time.Time

	rpcProcessor  thrift.TProcessor
	rpcServer     thrift.TServer
	rpcThreadPool *rpcThreadPool

	stats    *engineStats
	pid      int
	hostname string
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	this.stats = newEngineStats()

	return
}
