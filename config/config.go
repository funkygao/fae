package config

import (
	conf "github.com/funkygao/jsconf"
)

var (
	Servants *ConfigServant
)

type ConfigServant struct {
	WatchdogInterval    int
	ProfilerMaxBodySize int
	ProfilerRate        int

	// distribute load accross servers
	PeersCooperate bool

	Mongodb  *ConfigMongodb
	Memcache *ConfigMemcache
	Lcache   *ConfigLcache
}

func init() {
	Servants = new(ConfigServant)
}

func LoadServants(cf *conf.Conf) {
	Servants.WatchdogInterval = cf.Int("watchdog_interval", 60*10)
	Servants.PeersCooperate = cf.Bool("peers_cooperate", false)
	Servants.ProfilerMaxBodySize = cf.Int("profiler_max_body_size", 1<<10)
	Servants.ProfilerRate = cf.Int("profiler_rate", 1) // default 1/1000

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

	// lcache section
	Servants.Lcache = new(ConfigLcache)
	section, err = cf.Section("lcache")
	if err != nil {
		panic(err)
	}
	Servants.Lcache.loadConfig(section)
}
