package config

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigMongodbServer struct {
	Pool         string
	Host         string
	Port         string
	User         string
	Pass         string
	DbName       string
	ReplicaSet   string
	ShardBaseNum int

	uri string // cache of op result
}

func (this *ConfigMongodbServer) loadConfig(section *conf.Conf) {
	this.Pool = section.String("pool", "")
	this.Host = section.String("host", "")
	this.Port = section.String("port", "27017")
	this.DbName = section.String("db", "")
	this.ShardBaseNum = section.Int("shard_base_num", this.ShardBaseNum)
	this.User = section.String("user", "")
	this.Pass = section.String("pass", "")
	this.ReplicaSet = section.String("replicaSet", "")
	if this.Host == "" ||
		this.Port == "" ||
		this.Pool == "" ||
		this.DbName == "" {
		panic("required field missing")
	}

	// http://docs.mongodb.org/manual/reference/connection-string/
	this.uri = "mongodb://" + this.Host + ":" + this.Port + "/"
	if this.DbName != "" {
		this.uri += this.DbName + "/"
	}
	if this.ReplicaSet != "" {
		this.uri += "?replicaSet=" + this.ReplicaSet
	}

}

func (this *ConfigMongodbServer) Uri() string {
	return this.uri
}

type ConfigMongodb struct {
	DebugProtocol         bool
	DebugHeartbeat        bool
	ShardBaseNum          int
	ShardStrategy         string
	ConnectTimeout        time.Duration
	IoTimeout             time.Duration
	MaxIdleConnsPerServer int
	MaxConnsPerServer     int
	HeartbeatInterval     int
	Breaker               ConfigBreaker
	Servers               map[string]*ConfigMongodbServer // key is pool
}

func (this *ConfigMongodb) Enabled() bool {
	return len(this.Servers) > 0
}

func (this *ConfigMongodb) loadConfig(cf *conf.Conf) {
	this.ShardBaseNum = cf.Int("shard_base_num", 100000)
	this.DebugProtocol = cf.Bool("debug_protocol", false)
	this.DebugHeartbeat = cf.Bool("debug_heartbeat", false)
	this.ShardStrategy = cf.String("shard_strategy", "legacy")
	this.ConnectTimeout = cf.Duration("connect_timeout", 4*time.Second)
	this.IoTimeout = cf.Duration("io_timeout", 30*time.Second)
	this.MaxIdleConnsPerServer = cf.Int("max_idle_conns_per_server", 2)
	this.MaxConnsPerServer = cf.Int("max_conns_per_server",
		this.MaxIdleConnsPerServer*5)
	this.HeartbeatInterval = cf.Int("heartbeat_interval", 120)
	section, err := cf.Section("breaker")
	if err == nil {
		this.Breaker.loadConfig(section)
	}
	this.Servers = make(map[string]*ConfigMongodbServer)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMongodbServer)
		server.ShardBaseNum = this.ShardBaseNum
		server.loadConfig(section)
		this.Servers[server.Pool] = server
	}

	log.Debug("mongodb conf: %+v", *this)
}
