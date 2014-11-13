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

	// distribute load accross servers
	PeersReplica          int
	PeerGroupAddr         string
	PeerHeartbeatInterval int
	PeerDeadThreshold     float64

	Mongodb  *ConfigMongodb
	Memcache *ConfigMemcache
	Lcache   *ConfigLcache
	Proxy    *ConfigProxy
	Mysql    *ConfigMysql
	Redis    *ConfigRedis // TODO
}

func LoadServants(cf *conf.Conf) {
	Servants.DataCenterId = cf.Int("data_center_id", 1)
	Servants.AgentId = cf.Int("agent_id", 1)
	Servants.CallSlowThreshold = cf.Duration("call_slow_threshold", 2*time.Second)
	Servants.StatsOutputInterval = cf.Duration("stats_output_interval", 10*time.Minute)
	Servants.ProfilerMaxBodySize = cf.Int("profiler_max_body_size", 1<<10)
	Servants.ProfilerRate = cf.Int("profiler_rate", 1) // default 1/1000
	Servants.PeersReplica = cf.Int("peer_replicas", 3)
	Servants.PeerHeartbeatInterval = cf.Int("peer_heartbeat_interval", 0)
	Servants.PeerDeadThreshold = cf.Float("peer_dead_threshold",
		float64(Servants.PeerHeartbeatInterval)*3)
	Servants.PeerGroupAddr = cf.String("peer_group_addr", "224.0.0.2:19850")

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

	// proxy section
	Servants.Proxy = new(ConfigProxy)
	section, err = cf.Section("proxy")
	if err == nil {
		Servants.Proxy.loadConfig(section)
	}

	log.Debug("servants: %+v", *Servants)
}

func (this *ConfigServant) PeerEnabled() bool {
	return this.PeerHeartbeatInterval > 0
}
