package engine

import (
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

	if err := this.peer.Start(); err != nil {
		log.Error(err)
	}

	<-this.launchRpcServe()

	log.Info("Engine terminated")
}
