package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant"
	"time"
)

type Engine struct {
	conf *engineConfig

	StartedAt time.Time

	svt           *servant.FunServantImpl
	rpcProcessor  thrift.TProcessor
	rpcServer     thrift.TServer
	rpcThreadPool *rpcThreadPool

	stats    *engineStats
	pid      int
	hostname string

	stopChan chan bool
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	this.conf = new(engineConfig)
	this.stats = newEngineStats()
	this.stopChan = make(chan bool)

	return
}
