package config

import (
	"github.com/funkygao/golib/ip"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"strings"
)

// the root of config, which will load sections one by one
type ConfigEngine struct {
	*conf.Conf

	EtcdServers  []string
	EtcdSelfAddr string

	HttpListenAddr  string
	PprofListenAddr string
	MetricsLogfile  string

	Rpc      *ConfigRpc
	Servants *ConfigServant
}

func LoadEngineConfig(cf *conf.Conf) {
	Engine = new(ConfigEngine)
	Engine.LoadConfig(cf)
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
