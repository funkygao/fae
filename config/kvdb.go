package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigKvdb struct {
	path       string
	servletNum int
	enabled    bool
}

func (this *ConfigKvdb) loadConfig(cf *conf.Conf) {
	this.path = cf.String("path", "/tmp/kvdb")
	this.servletNum = cf.Int("servlet_num", 0)
	this.enabled = true

	log.Debug("kvdb: %+v", *this)
}

func (this *ConfigKvdb) Enabled() bool {
	return this.enabled
}
