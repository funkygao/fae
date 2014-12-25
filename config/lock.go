package config

import (
	conf "github.com/funkygao/jsconf"
	"time"
)

type ConfigLock struct {
	MaxItems int
	Expires  time.Duration
}

func (this *ConfigLock) LoadConfig(cf *conf.Conf) {
	this.MaxItems = cf.Int("max_items", 1<<20)
	this.Expires = cf.Duration("expires", 10*time.Second)
}
