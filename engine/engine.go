package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/peer"
	"time"
)

type Engine struct {
	conf *engineConfig

	configFile string
	StartedAt  time.Time

	rpcProcessor thrift.TProcessor
	rpcServer    thrift.TServer

	peer *peer.Peer

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
