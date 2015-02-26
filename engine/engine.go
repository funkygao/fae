package engine

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant"
	conf "github.com/funkygao/jsconf"
	"github.com/funkygao/thrift/lib/go/thrift"
	"time"
)

type Engine struct {
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
	this.stats = newEngineStats()
	this.stopChan = make(chan bool)

	return
}

func (this *Engine) LoadConfig(configFile string, cf *conf.Conf) *Engine {
	config.LoadEngineConfig(configFile, cf)
	return this
}
