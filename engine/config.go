package engine

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	conf "github.com/funkygao/jsconf"
	"time"
)

type configRpc struct {
	clientSlowThreshold float64 // in seconds per connection
	callSlowThreshold   float64 // in seconds per call
	listenAddr          string
	clientTimeout       time.Duration
	framed              bool
	protocol            string
}

func (this *configRpc) loadConfig(section *conf.Conf) {
	this.listenAddr = section.String("listen_addr", "")
	if this.listenAddr == "" {
		panic("Empty listen_addr")
	}

	this.clientSlowThreshold = section.Float("client_slow_threshold", 5)
	this.callSlowThreshold = section.Float("call_slow_threshold", 5)
	this.clientTimeout = time.Duration(section.Int("client_timeout", 0)) * time.Second
	this.framed = section.Bool("framed", false)
	this.protocol = section.String("protocol", "binary")

	log.Debug("rpc: %+v", *this)
}

type engineConfig struct {
	*conf.Conf

	httpListenAddr        string
	peerGroupAddr         string
	peerHeartbeatInterval int

	rpc *configRpc
}

func (this *Engine) LoadConfigFile() *Engine {
	log.Debug("Loading config file %s", this.configFile)

	cf := new(engineConfig)
	var err error
	cf.Conf, err = conf.Load(this.configFile)
	if err != nil {
		panic(err)
	}

	this.conf = cf
	this.doLoadConfig()

	// RegisterHttpApi is ready
	this.setupHttpServ()

	// when config loaded, create the servants
	this.rpcProcessor = rpc.NewFunServantProcessor(servant.NewFunServant(config.Servants))

	this.peer = newPeer(this.conf.peerGroupAddr, this.conf.peerHeartbeatInterval)

	return this
}

func (this *Engine) doLoadConfig() {
	this.conf.httpListenAddr = this.conf.String("http_listen_addr", "")
	this.conf.peerHeartbeatInterval = this.conf.Int("peer_heartbeat_interval", 30)
	this.conf.peerGroupAddr = this.conf.String("peer_group_addr", "224.0.0.2:19850")

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
}
