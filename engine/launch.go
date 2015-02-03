package engine

import (
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
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
	cpuNumConfig := config.Engine.String("cpu_num", "auto")
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
	go this.stats.Start(this.StartedAt,
		config.Engine.Rpc.StatsOutputInterval,
		config.Engine.MetricsLogfile)

	this.launchHttpServ()
	defer this.stopHttpServ()

	select {
	case <-this.launchRpcServe():
	case <-this.stopChan:
	}

	log.Info("Engine terminated")
}

func (this *Engine) UnregisterEtcd() {
	if config.Engine.IsProxyOnly() {
		return
	}

	if config.Engine.EtcdSelfAddr != "" {
		if !etclib.IsConnected() {
			if err := etclib.Dial(config.Engine.EtcdServers); err != nil {
				log.Error("etcd[%+v]: %s", config.Engine.EtcdServers, err)
				return
			}
		}

		etclib.ShutdownService(config.Engine.EtcdSelfAddr, etclib.SERVICE_FAE)
		etclib.Close()

		log.Info("etcd self[%s] unregistered", config.Engine.EtcdSelfAddr)
	}
}
