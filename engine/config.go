package engine

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/ip"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"strings"
	"time"
)

type engineConfig struct {
	*conf.Conf

	EtcdServers  []string
	EtcdSelfAddr string

	httpListenAddr  string
	pprofListenAddr string
	metricsLogfile  string

	rpc *configRpc
}

func (this *Engine) LoadConfig(cf *conf.Conf) *Engine {
	this.conf.Conf = cf

	this.conf.EtcdServers = cf.StringList("etcd_servers", nil)
	if len(this.conf.EtcdServers) > 0 {
		this.conf.EtcdSelfAddr = cf.String("etcd_self_addr", "")
		if strings.HasPrefix(this.conf.EtcdSelfAddr, ":") {
			// automatically get local ip addr
			myIp := ip.LocalIpv4Addrs()[0]
			this.conf.EtcdSelfAddr = myIp + this.conf.EtcdSelfAddr
		}
	}
	this.conf.httpListenAddr = this.conf.String("http_listen_addr", "")
	this.conf.pprofListenAddr = this.conf.String("pprof_listen_addr", "")
	this.conf.metricsLogfile = this.conf.String("metrics_logfile", "metrics.log")

	// rpc section
	this.conf.rpc = new(configRpc)
	section, err := this.conf.Section("rpc")
	if err != nil {
		panic(err)
	}
	this.conf.rpc.loadConfig(section)

	section, err = this.conf.Section("servants")
	if err != nil {
		panic(err)
	}
	config.LoadServants(section)

	log.Debug("engine conf: %+v", *this.conf)

	return this
}

type configRpc struct {
	listenAddr             string
	sessionSlowThreshold   time.Duration // per session
	sessionTimeout         time.Duration
	ioTimeout              time.Duration
	bufferSize             int // network IO read/write buffer
	framed                 bool
	protocol               string
	statsOutputInterval    time.Duration
	maxOutstandingSessions int
}

func (this *configRpc) loadConfig(section *conf.Conf) {
	this.listenAddr = section.String("listen_addr", "")
	if this.listenAddr == "" {
		panic("Empty listen_addr")
	}

	this.sessionSlowThreshold = section.Duration("session_slow_threshold", 10*time.Second)
	this.sessionTimeout = section.Duration("session_timeout", 30*time.Second)
	this.ioTimeout = section.Duration("io_timeout", 2*time.Second)
	this.statsOutputInterval = section.Duration("stats_output_interval", 10*time.Second)
	this.framed = section.Bool("framed", false)
	this.bufferSize = section.Int("buffer_size", 4<<10)
	this.protocol = section.String("protocol", "binary")
	this.maxOutstandingSessions = section.Int("max_outstanding_sessions", 20000)

	log.Debug("rpc conf: %+v", *this)
}
