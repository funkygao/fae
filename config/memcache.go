package config

import (
	log "code.google.com/p/log4go"
	"fmt"
	conf "github.com/funkygao/jsconf"
)

type ConfigMemcacheServer struct {
	host string
	hort string
}

func (this *ConfigMemcacheServer) loadConfig(section *conf.Conf) {
	this.host = section.String("host", "")
	if this.host == "" {
		panic("Empty memcache server host")
	}
	this.hort = section.String("port", "")
	if this.hort == "" {
		panic("Empty memcache server port")
	}

	log.Debug("memcache server: %+v", *this)
}

func (this *ConfigMemcacheServer) Address() string {
	return this.host + ":" + this.hort
}

type ConfigMemcache struct {
	HashStrategy string
	HashFunction string

	Servers map[string]*ConfigMemcacheServer // key is host:port(addr)
}

func (this *ConfigMemcache) loadConfig(cf *conf.Conf) {
	this.Servers = make(map[string]*ConfigMemcacheServer)
	this.HashStrategy = cf.String("hash_strategy", "standard")
	this.HashFunction = cf.String("hash_function", "crc32")
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMemcacheServer)
		server.loadConfig(section)
		this.Servers[server.Address()] = server
	}

	log.Debug("memcache: %+v", *this)
}
