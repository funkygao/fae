package config

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigMemcacheServer struct {
	pool string
	host string
	port string
}

func (this *ConfigMemcacheServer) loadConfig(section *conf.Conf) {
	this.host = section.String("host", "")
	if this.host == "" {
		panic("Empty memcache server host")
	}
	this.port = section.String("port", "")
	if this.port == "" {
		panic("Empty memcache server port")
	}
	this.pool = section.String("pool", "default")

	log.Debug("memcache server: %+v", *this)
}

func (this *ConfigMemcacheServer) Address() string {
	return this.host + ":" + this.port
}

type ConfigMemcache struct {
	HashStrategy          string
	Timeout               int
	MaxIdleConnsPerServer int
	enabled               bool

	Servers map[string]*ConfigMemcacheServer // key is host:port(addr)
}

func (this *ConfigMemcache) ServerList() []string {
	servers := make([]string, len(this.Servers))
	i := 0
	for addr, _ := range this.Servers {
		servers[i] = addr
		i += 1
	}

	return servers
}

func (this *ConfigMemcache) Enabled() bool {
	return this.enabled
}

func (this *ConfigMemcache) loadConfig(cf *conf.Conf) {
	this.Servers = make(map[string]*ConfigMemcacheServer)
	this.HashStrategy = cf.String("hash_strategy", "standard")
	this.Timeout = cf.Int("timeout", 4)
	this.MaxIdleConnsPerServer = cf.Int("max_idle_conns_per_server", 3)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMemcacheServer)
		server.loadConfig(section)
		this.Servers[server.Address()] = server
	}
	this.enabled = true

	log.Debug("memcache: %+v", *this)
}
