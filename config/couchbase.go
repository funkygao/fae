package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigCouchbase struct {
	Server string // TODO cluster
}

func (this *ConfigCouchbase) loadConfig(cf *conf.Conf) {
	this.Server = cf.String("server", "http://localhost:8091/")
	log.Debug("couchbase: %+v", *this)
}
