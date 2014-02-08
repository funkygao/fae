package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigLcache struct {
	LruMaxItems int
}

func (this *ConfigLcache) loadConfig(cf *conf.Conf) {
	this.LruMaxItems = cf.Int("lru_max_items", 1<<30)

	log.Debug("lcache: %+v", *this)
}
