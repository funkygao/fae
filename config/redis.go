package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigRedisServer struct {
	Host        string
	Port        string
	MaxIdle     int
	IdleTimeout time.Duration
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

	log.Debug("redis conf: %+v", *this)
}
