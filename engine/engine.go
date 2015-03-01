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

	svt          *servant.FunServantImplWrapper
	rpcProcessor thrift.TProcessor
	rpcServer    thrift.TServer

	pid      int
	hostname string

	stopChan chan bool
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	this.stopChan = make(chan bool)

	return
}

func (this *Engine) LoadConfig(configFile string, cf *conf.Conf) *Engine {
	config.LoadEngineConfig(configFile, cf)
	return this
}
