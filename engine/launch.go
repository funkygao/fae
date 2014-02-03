package engine

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/golib/signal"
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

	<-this.launchRpcServe()

	log.Info("Engine terminated")
}
