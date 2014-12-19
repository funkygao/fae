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

	signal.RegisterSignalHandler(syscall.SIGUSR1, func(sig os.Signal) {
		// graceful shutdown
		this.StopRpcServe()
	})

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

	// start the stats counter
	go this.stats.Start(this.StartedAt, this.conf.rpc.statsOutputInterval,
		this.conf.metricsLogfile)

	this.launchHttpServ()
	defer this.stopHttpServ()

	select {
	case <-this.launchRpcServe():
	case <-this.stopChan:
	}

	log.Info("Engine terminated")
}

func (this *Engine) UnregisterEtcd() {
	if this.conf.EtcdSelfAddr != "" {
		if err := etclib.Dial(this.conf.EtcdServers); err != nil {
			log.Error("etcd[%+v]: %s", this.conf.EtcdServers, err)
			return
		}

		etclib.ShutdownService(this.conf.EtcdSelfAddr, etclib.SERVICE_FAE)
		etclib.Close()

		log.Info("etcd self[%s] unregistered", this.conf.EtcdSelfAddr)
	}
}
