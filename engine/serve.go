package engine

import (
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

	this.launchHttpServ()
	defer this.stopHttpServ()

	<-this.launchRpcServe()
}
