package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

var (
	Servants *ConfigServant
)

func init() {
	Servants = new(ConfigServant)
}

type ConfigServant struct {
	DataCenterId int
	AgentId      int

	CallSlowThreshold   time.Duration
	StatsOutputInterval time.Duration
	ProfilerMaxBodySize int
	ProfilerRate        int
	SessionEntries      int // LRU cache volumn

	Mongodb   *ConfigMongodb
	Memcache  *ConfigMemcache
	Lcache    *ConfigLcache
	Proxy     *ConfigProxy
	Mysql     *ConfigMysql
	Redis     *ConfigRedis // TODO
	Couchbase *ConfigCouchbase
}

func LoadServants(cf *conf.Conf) {
	Servants.DataCenterId = cf.Int("data_center_id", 1)
	Servants.AgentId = cf.Int("agent_id", 1)
	Servants.SessionEntries = cf.Int("session_entries", 20<<10)
	Servants.CallSlowThreshold = cf.Duration("call_slow_threshold", 2*time.Second)
	Servants.StatsOutputInterval = cf.Duration("stats_output_interval", 10*time.Minute)
	Servants.ProfilerMaxBodySize = cf.Int("profiler_max_body_size", 1<<10)
	Servants.ProfilerRate = cf.Int("profiler_rate", 1) // default 1/1000

	// mongodb section
	Servants.Mongodb = new(ConfigMongodb)
	section, err := cf.Section("mongodb")
	if err == nil {
		Servants.Mongodb.loadConfig(section)
	}

	Servants.Mysql = new(ConfigMysql)
	section, err = cf.Section("mysql")
	if err == nil {
		Servants.Mysql.loadConfig(section)
	}

	Servants.Redis = new(ConfigRedis)
	section, err = cf.Section("redis")
	if err == nil {
		Servants.Redis.loadConfig(section)
	}

	// memcached section
	Servants.Memcache = new(ConfigMemcache)
	section, err = cf.Section("memcache")
	if err == nil {
		Servants.Memcache.loadConfig(section)
	}

	// lcache section
	Servants.Lcache = new(ConfigLcache)
	section, err = cf.Section("lcache")
	if err == nil {
		Servants.Lcache.loadConfig(section)
	}

	// couchbase section
	Servants.Couchbase = new(ConfigCouchbase)
	section, err = cf.Section("couchbase")
	if err == nil {
		Servants.Couchbase.loadConfig(section)
	}

	// proxy section
	Servants.Proxy = new(ConfigProxy)
	section, err = cf.Section("proxy")
	if err == nil {
		Servants.Proxy.loadConfig(section)
	}

	log.Debug("servants conf: %+v", *Servants)
}
