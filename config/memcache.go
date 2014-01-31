package config

import (
	log "code.google.com/p/log4go"
	"fmt"
	conf "github.com/funkygao/jsconf"
)

type ConfigMemcacheServer struct {
	Host string
	Port string
}

func (this *ConfigMemcacheServer) loadConfig(section *conf.Conf) {
	this.Host = section.String("host", "")
	if this.Host == "" {
		panic("Empty memcache server host")
	}
	this.Port = section.String("port", "")
	if this.Port == "" {
		panic("Empty memcache server port")
	}

	log.Debug("memcache server: %+v", *this)
}

type ConfigMemcache struct {
	HashStrategy string
	HashFunction string
	Servers      map[string]*ConfigMemcacheServer // key is host:port
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
		this.Servers[server.Host+":"+server.Port] = server
	}

	log.Debug("memcache: %+v", *this)
}
