package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigLock struct {
	MaxItems int
	Expires  time.Duration
	enabled  bool
}

func (this *ConfigLock) LoadConfig(cf *conf.Conf) {
	this.MaxItems = cf.Int("max_items", 1<<20)
	this.Expires = cf.Duration("expires", time.Second*10)

	this.enabled = true

	log.Debug("lock conf: %+v", *this)
}

func (this *ConfigLock) Enabled() bool {
	return this.enabled && this.MaxItems > 0
}
