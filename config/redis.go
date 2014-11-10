package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigRedisServer struct {
}

type ConfigRedis struct {
	Breaker ConfigBreaker
	Servers map[string]*ConfigRedisServer
}

func (this *ConfigRedis) loadConfig(cf *conf.Conf) {
	section, err := cf.Section("breaker")
	if err == nil {
		this.Breaker.loadConfig(section)
	}

	log.Debug("redis: %+v", *this)
}
