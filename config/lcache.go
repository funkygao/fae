package config

import (
	log "code.google.com/p/log4go"
	conf "github.com/funkygao/jsconf"
)

type ConfigLcache struct {
	LruMaxItems int
}

func (this *ConfigLcache) loadConfig(cf *conf.Conf) {
	this.LruMaxItems = cf.Int("lru_max_items", 1<<30)

	log.Debug("lcache: %+v", *this)
}
