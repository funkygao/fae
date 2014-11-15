package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigCouchbase struct {
}

func (this *ConfigCouchbase) loadConfig(cf *conf.Conf) {

	log.Debug("couchbase: %+v", *this)
}
