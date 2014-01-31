package config

import (
	conf "github.com/funkygao/jsconf"
)

var (
	servantConfig *ConfigServant
)

type ConfigServant struct {
	mongodb  *ConfigMongodb
	memcache *ConfigMemcache
}

func init() {
	servantConfig = new(ConfigServant)
}

func LoadServants(cf *conf.Conf) {
	// mongodb section
	servantConfig.mongodb = new(ConfigMongodb)
	section, err := cf.Section("mongodb")
	if err != nil {
		panic(err)
	}
	servantConfig.mongodb.loadConfig(section)

	// memcached section
	servantConfig.memcache = new(ConfigMemcache)
	section, err = cf.Section("memcache")
	if err != nil {
		panic(err)
	}
	servantConfig.memcache.loadConfig(section)
}
