package config

import (
	"github.com/funkygao/golib/ip"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"os"
	"strings"
	"time"
)

// the root of config, which will load sections one by one
type ConfigEngine struct {
	*conf.Conf

	configFile         string
	configFileLastStat os.FileInfo

	EtcdServers  []string
	EtcdSelfAddr string

	HttpListenAddr  string
	PprofListenAddr string
	MetricsLogfile  string

	ReloadWatchdogInterval time.Duration
	ServerMode             bool // if false then client only mode

	Rpc      *ConfigRpc
	Servants *ConfigServant

	ReloadedChan chan ConfigEngine
}

func (this *ConfigEngine) LoadConfig(cf *conf.Conf) {
	this.Conf = cf

	this.HttpListenAddr = this.String("http_listen_addr", "")
	this.PprofListenAddr = this.String("pprof_listen_addr", "")
	this.MetricsLogfile = this.String("metrics_logfile", "metrics.log")
	this.ReloadWatchdogInterval = this.Duration("reload_watchdog_interval", time.Second)
	this.ServerMode = this.Bool("server_mode", true)

	// rpc section
	this.Rpc = new(ConfigRpc)
	section, err := this.Section("rpc")
	if err != nil {
		panic(err)
	}
	this.Rpc.LoadConfig(section)

	// servants section
	this.Servants = new(ConfigServant)
	section, err = this.Section("servants")
	if err != nil {
		panic(err)
	}
	this.Servants.LoadConfig(this.Rpc.ListenAddr, section)

	// after load all configs, calculate EtcdSelfAddr
	this.EtcdServers = cf.StringList("etcd_servers", nil)
	if len(this.EtcdServers) > 0 {
		this.EtcdSelfAddr = this.Rpc.ListenAddr
		if strings.HasPrefix(this.EtcdSelfAddr, ":") {
			// automatically get local ip addr
			this.EtcdSelfAddr = ip.LocalIpv4Addrs()[0] + this.EtcdSelfAddr
		}
	}

	log.Debug("engine conf: %+v", *this)
}

func (this *ConfigEngine) runWatchdog() {
	ticker := time.NewTicker(this.ReloadWatchdogInterval)
	defer ticker.Stop()

	for _ = range ticker.C {
		stat, _ := os.Stat(this.configFile)
		if stat.ModTime() != this.configFileLastStat.ModTime() {
			this.configFileLastStat = stat

			cf, err := conf.Load(this.configFile)
			if err != nil {
				panic(err)
			}

			log.Info("config[%s] reloaded", this.configFile)
			this.LoadConfig(cf)
			this.ReloadedChan <- *this
		}
	}

}
