package engine

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/peer"
	"github.com/funkygao/fae/servant"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/signal"
	log "github.com/funkygao/log4go"
	"os"
	"syscall"
	"time"
)

func (this *Engine) ServeForever() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()

	signal.IgnoreSignal(syscall.SIGHUP)

	log.Info("Launching Engine...")

	// start the stats counter
	this.stats.Start(this.StartedAt)

	this.launchHttpServ()
	defer this.stopHttpServ()

	// when config loaded, create the servants
	svr := servant.NewFunServant(config.Servants)
	this.rpcProcessor = rpc.NewFunServantProcessor(svr)
	svr.Start()

	this.peer = peer.NewPeer(this.conf.peerGroupAddr,
		this.conf.peerHeartbeatInterval, this.conf.peerDeadThreshold)
	if err := this.peer.Start(); err != nil {
		log.Error(err)
	}

	<-this.launchRpcServe()

	log.Info("Engine terminated")
}
