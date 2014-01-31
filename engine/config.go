package engine

import (
	log "code.google.com/p/log4go"
	"fmt"
	conf "github.com/daviddengcn/go-ljson-conf"
)

type ConfigMemcache struct {
	host string
	port int
}

func (this *ConfigMemcache) loadConfig(section *conf.Conf) {
	this.host = section.String("host", "")
	this.port = section.Int("port", 0)
	if this.host == "" || this.port == 0 {
		panic("required filed")
	}
}

type ConfigMongodb struct {
	host       string
	port       int
	user, pass string
	db         string
	replicaSet string
}

func (this *ConfigMongodb) loadConfig(section *conf.Conf) {
	this.host = section.String("host", "")
	this.port = section.Int("port", 0)
	this.user = section.String("user", "")
	this.pass = section.String("pass", "")
	this.replicaSet = section.String("replicaSet", "")
	if this.host == "" ||
		this.port == 0 ||
		this.replicaSet == "" {
		panic("required filed")
	}
}

type Config struct {
	*conf.Conf

	rpcListenAddr  string
	httpListenAddr string

	mongos    []*ConfigMongodb
	memcaches []*ConfigMemcache
}

func (this *Engine) LoadConfigFile(fn string) *Engine {
	log.Debug("Loading config file %s", fn)

	config := new(Config)
	var err error
	config.Conf, err = conf.Load(fn)
	if err != nil {
		panic(err)
	}

	this.conf = config
	this.doLoadConfig()

	return this
}

func (this *Engine) doLoadConfig() {
	this.conf.httpListenAddr = this.conf.String("http_listen_addr", "")
	this.conf.rpcListenAddr = this.conf.String("rpc_listen_addr", "")
	if this.conf.rpcListenAddr == "" {
		panic("rpc_listen_addr empty")
	}

	this.conf.mongos = make([]*ConfigMongodb, 0)
	this.conf.memcaches = make([]*ConfigMemcache, 0)

	for i := 0; i < len(this.conf.List("mongodb", nil)); i++ {
		section, err := this.conf.Section(fmt.Sprintf("mongodb[%d]", i))
		if err != nil {
			panic(err)
		}

		cf := new(ConfigMongodb)
		cf.loadConfig(section)
		this.conf.mongos = append(this.conf.mongos, cf)
	}

	for i := 0; i < len(this.conf.List("memcached", nil)); i++ {
		section, err := this.conf.Section(fmt.Sprintf("memcached[%d]", i))
		if err != nil {
			panic(err)
		}

		cf := new(ConfigMemcache)
		cf.loadConfig(section)
		this.conf.memcaches = append(this.conf.memcaches, cf)
	}

	log.Debug("config %+v", *this.conf)
}
