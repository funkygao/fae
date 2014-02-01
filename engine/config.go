package engine

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fxi/config"
	"github.com/funkygao/fxi/servant"
	"github.com/funkygao/fxi/servant/gen-go/fun/rpc"
	conf "github.com/funkygao/jsconf"
)

type configRpc struct {
	listenAddr string
	framed     bool
	protocol   string
}

func (this *configRpc) loadConfig(section *conf.Conf) {
	this.listenAddr = section.String("listen_addr", "")
	if this.listenAddr == "" {
		panic("Empty listen_addr")
	}

	this.framed = section.Bool("framed", false)
	this.protocol = section.String("protocol", "binary")

	log.Debug("rpc: %+v", *this)
}

type engineConfig struct {
	*conf.Conf

	httpListenAddr string

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

	// when config loaded, create the servants
	this.rpcProcessor = rpc.NewFunServantProcessor(servant.NewFunServant(config.Servants))

	return this
}

func (this *Engine) doLoadConfig() {
	this.conf.httpListenAddr = this.conf.String("http_listen_addr", "")

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
