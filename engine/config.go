package engine

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/peer"
	"github.com/funkygao/fae/servant"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type configRpc struct {
	sessionSlowThreshold float64 // in seconds per session
	callSlowThreshold    float64 // in seconds per call
	listenAddr           string
	clientTimeout        time.Duration
	framed               bool
	protocol             string
	debugSession         bool
	tcpNoDelay           bool
}

func (this *configRpc) loadConfig(section *conf.Conf) {
	this.listenAddr = section.String("listen_addr", "")
	if this.listenAddr == "" {
		panic("Empty listen_addr")
	}

	this.sessionSlowThreshold = section.Float("session_slow_threshold", 5)
	this.callSlowThreshold = section.Float("call_slow_threshold", 5)
	this.clientTimeout = time.Duration(section.Int("client_timeout", 0)) * time.Second
	this.framed = section.Bool("framed", false)
	this.protocol = section.String("protocol", "binary")
	this.tcpNoDelay = section.Bool("tcp_nodelay", true)
	this.debugSession = section.Bool("debug_session", false)

	log.Debug("rpc: %+v", *this)
}

type engineConfig struct {
	*conf.Conf

	httpListenAddr        string
	peerGroupAddr         string
	peerHeartbeatInterval int
	peerDeadThreshold     float64

	rpc *configRpc
}

func (this *Engine) LoadConfigFile() *Engine {
	log.Info("Engine[%s] loading config file %s", BuildID, this.configFile)

	cf := new(engineConfig)
	var err error
	cf.Conf, err = conf.Load(this.configFile)
	if err != nil {
		panic(err)
	}

	this.conf = cf
	this.doLoadConfig()

	// when config loaded, create the servants
	svr := servant.NewFunServant(config.Servants)
	this.rpcProcessor = rpc.NewFunServantProcessor(svr)
	svr.Start()

	this.peer = peer.NewPeer(this.conf.peerGroupAddr,
		this.conf.peerHeartbeatInterval, this.conf.peerDeadThreshold)

	return this
}

func (this *Engine) doLoadConfig() {
	this.conf.httpListenAddr = this.conf.String("http_listen_addr", "")
	this.conf.peerHeartbeatInterval = this.conf.Int("peer_heartbeat_interval", 30)
	this.conf.peerDeadThreshold = this.conf.Float("peer_dead_threshold", 30)
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

	log.Debug("engine: %+v", *this.conf)
}
