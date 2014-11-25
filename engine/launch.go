package engine

import (
	"github.com/funkygao/etclib"
	"github.com/funkygao/golib/signal"
	log "github.com/funkygao/log4go"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func (this *Engine) ServeForever() {
	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()

	signal.IgnoreSignal(syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGSTOP)

	var (
		totalCpus int
		maxProcs  int
	)
	totalCpus = runtime.NumCPU()
	cpuNumConfig := this.conf.String("cpu_num", "auto")
	if cpuNumConfig == "auto" {
		maxProcs = totalCpus/2 + 1
	} else if cpuNumConfig == "max" {
		maxProcs = totalCpus
	} else {
		maxProcs, _ = strconv.Atoi(cpuNumConfig)
	}
	runtime.GOMAXPROCS(maxProcs)
	log.Info("Launching Engine with %d/%d CPUs...", maxProcs, totalCpus)

	if this.conf.EtcdSelfAddr != "" {
		etclib.Init(this.conf.EtcdServers, "dw")
		etclib.BootFae(this.conf.EtcdSelfAddr)
	}

	// start the stats counter
	go this.stats.Start(this.StartedAt, this.conf.rpc.statsOutputInterval,
		this.conf.metricsLogfile)

	this.launchHttpServ()
	defer this.stopHttpServ()

	<-this.launchRpcServe()

	log.Info("Engine terminated")
}

func (this *Engine) Stop() {
	if this.conf.EtcdSelfAddr != "" {
		etclib.Init(this.conf.EtcdServers, "dw")
		etclib.ShutdownFae(this.conf.EtcdSelfAddr)
	}
}
