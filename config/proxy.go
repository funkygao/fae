package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigProxy struct {
	PoolCapacity int
	IdleTimeout  time.Duration
	enabled      bool
}

func (this *ConfigProxy) loadConfig(cf *conf.Conf) {
	this.PoolCapacity = cf.Int("pool_capacity", 10)
	this.IdleTimeout = time.Duration(cf.Int("idle_timeout", 600)) * time.Second
	this.enabled = true

	log.Debug("proxy: %+v", *this)
}

func (this *ConfigProxy) Enabled() bool {
	return this.enabled
}
