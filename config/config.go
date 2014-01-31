package config

import (
	conf "github.com/funkygao/jsconf"
)

var (
	Servants *ConfigServant
)

type ConfigServant struct {
	Mongodb  *ConfigMongodb
	Memcache *ConfigMemcache
}

func init() {
	Servants = new(ConfigServant)
}

func LoadServants(cf *conf.Conf) {
	// mongodb section
	Servants.Mongodb = new(ConfigMongodb)
	section, err := cf.Section("mongodb")
	if err != nil {
		panic(err)
	}
	Servants.Mongodb.loadConfig(section)

	// memcached section
	Servants.Memcache = new(ConfigMemcache)
	section, err = cf.Section("memcache")
	if err != nil {
		panic(err)
	}
	Servants.Memcache.loadConfig(section)
}
