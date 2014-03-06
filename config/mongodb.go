package config

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
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

	log.Debug("mongodb server: %+v", *this)
}

// http://docs.mongodb.org/manual/reference/connection-string/
func (this *ConfigMongodbServer) Address() string {
	addr := "mongodb://" + this.Host + ":" + this.Port + "/"
	if this.DbName != "" {
		addr += this.DbName + "/"
	}
	if this.ReplicaSet != "" {
		addr += "?replicaSet=" + this.ReplicaSet
	}
	return addr
}

func (this *ConfigMongodbServer) Url() string {
	return this.Address()
}

type ConfigMongodb struct {
	DebugProtocol         bool
	DebugHeartbeat        bool
	ShardBaseNum          int
	ShardStrategy         string
	ConnectTimeout        int
	IoTimeout             int
	MaxIdleConnsPerServer int
	HeartbeatInterval     int
	Breaker               ConfigBreaker
	Servers               map[string]*ConfigMongodbServer // key is pool

	enabled bool
}

func (this *ConfigMongodb) Enabled() bool {
	return this.enabled
}

func (this *ConfigMongodb) loadConfig(cf *conf.Conf) {
	this.enabled = true
	this.ShardBaseNum = cf.Int("shard_base_num", 100000)
	this.DebugProtocol = cf.Bool("debug_protocol", false)
	this.DebugHeartbeat = cf.Bool("debug_heartbeat", false)
	this.ShardStrategy = cf.String("shard_strategy", "legacy")
	this.ConnectTimeout = cf.Int("connect_timeout", 4)
	this.IoTimeout = cf.Int("io_timeout", 30)
	this.MaxIdleConnsPerServer = cf.Int("max_idle_conns_per_server", 2)
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

	log.Debug("mongodb: %+v", *this)
}
