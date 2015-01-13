package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigServant struct {
	IdgenWorkerId int // TODO dynamic calculated

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
	Game      *ConfigGame
}

func (this *ConfigServant) LoadConfig(selfAddr string, cf *conf.Conf) {
	this.IdgenWorkerId = cf.Int("idgen_worker_id", 1)
	this.SessionEntries = cf.Int("session_entries", 20<<10)
	this.CallSlowThreshold = cf.Duration("call_slow_threshold", 2*time.Second)
	this.StatsOutputInterval = cf.Duration("stats_output_interval", 10*time.Minute)
	this.ProfilerMaxBodySize = cf.Int("profiler_max_body_size", 1<<10)
	this.ProfilerRate = cf.Int("profiler_rate", 1) // default 1/1000

	// mongodb section
	this.Mongodb = new(ConfigMongodb)
	section, err := cf.Section("mongodb")
	if err == nil {
		this.Mongodb.LoadConfig(section)
	}

	this.Mysql = new(ConfigMysql)
	section, err = cf.Section("mysql")
	if err == nil {
		this.Mysql.LoadConfig(section)
	}

	this.Redis = new(ConfigRedis)
	section, err = cf.Section("redis")
	if err == nil {
		this.Redis.LoadConfig(section)
	}

	// memcached section
	this.Memcache = new(ConfigMemcache)
	section, err = cf.Section("memcache")
	if err == nil {
		this.Memcache.LoadConfig(section)
	}

	// lcache section
	this.Lcache = new(ConfigLcache)
	section, err = cf.Section("lcache")
	if err == nil {
		this.Lcache.LoadConfig(section)
	}

	this.Game = new(ConfigGame)
	section, err = cf.Section("game")
	if err == nil {
		this.Game.LoadConfig(section)
	}

	// couchbase section
	this.Couchbase = new(ConfigCouchbase)
	section, err = cf.Section("couchbase")
	if err == nil {
		this.Couchbase.LoadConfig(section)
	}

	// proxy section
	this.Proxy = new(ConfigProxy)
	section, err = cf.Section("proxy")
	if err == nil {
		this.Proxy.LoadConfig(selfAddr, section)
	}

	log.Debug("servants conf: %+v", *this)
}
