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

	Rpc      *ConfigRpc
	Servants *ConfigServant
}

func (this *ConfigEngine) LoadConfig(cf *conf.Conf) {
	this.Conf = cf

	this.EtcdServers = cf.StringList("etcd_servers", nil)
	if len(this.EtcdServers) > 0 {
		this.EtcdSelfAddr = cf.String("etcd_self_addr", "")
		if strings.HasPrefix(this.EtcdSelfAddr, ":") {
			// automatically get local ip addr
			myIp := ip.LocalIpv4Addrs()[0]
			this.EtcdSelfAddr = myIp + this.EtcdSelfAddr
		}
	}
	this.HttpListenAddr = this.String("http_listen_addr", "")
	this.PprofListenAddr = this.String("pprof_listen_addr", "")
	this.MetricsLogfile = this.String("metrics_logfile", "metrics.log")
	this.ReloadWatchdogInterval = this.Duration("reload_watchdog_interval", time.Second)

	// rpc section
	this.Rpc = new(ConfigRpc)
	section, err := this.Section("rpc")
	if err != nil {
		panic(err)
	}
	this.Rpc.LoadConfig(section)

	this.Servants = new(ConfigServant)
	section, err = this.Section("servants")
	if err != nil {
		panic(err)
	}
	this.Servants.LoadConfig(section)

	log.Debug("engine conf: %+v", *this.Conf)
}
