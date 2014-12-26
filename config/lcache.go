package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigLcache struct {
	MaxItems int
}

func (this *ConfigLcache) LoadConfig(cf *conf.Conf) {
	this.MaxItems = cf.Int("lru_max_items", 1<<30)

	log.Debug("lcache conf: %+v", *this)
}

func (this *ConfigLcache) Enabled() bool {
	return this.MaxItems > 0
}
