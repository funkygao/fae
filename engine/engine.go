package engine

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/engine/plugin"
	_ "github.com/funkygao/fae/plugins"
	"github.com/funkygao/fae/servant"
	"github.com/funkygao/golib/null"
	conf "github.com/funkygao/jsconf"
	"github.com/funkygao/thrift/lib/go/thrift"
	"os"
	"time"
)

type Engine struct {
	StartedAt time.Time

	svt          *servant.FunServantImplWrapper
	rpcProcessor thrift.TProcessor
	rpcServer    thrift.TServer

	pid      int
	hostname string

	stopChan chan null.NullStruct
}

func NewEngine() *Engine {
	this := &Engine{stopChan: make(chan null.NullStruct),
		pid: os.Getpid()}
	this.hostname, _ = os.Hostname()
	return this
}

func (this *Engine) LoadConfig(cf *conf.Conf) *Engine {
	config.LoadEngineConfig(cf)

	if section, err := cf.Section("plugin"); err == nil {
		plugin.LoadPlugins(section)
	}
	return this
}
