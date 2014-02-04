package config

import (
	log "code.google.com/p/log4go"
	"fmt"
	conf "github.com/funkygao/jsconf"
)

type ConfigMongodbServer struct {
	Kind         string
	Host         string
	Port         string
	User         string
	Pass         string
	DbName       string
	ReplicaSet   string
	ShardBaseNum int
}

func (this *ConfigMongodbServer) loadConfig(section *conf.Conf) {
	this.Kind = section.String("kind", "")
	this.Host = section.String("host", "")
	this.Port = section.String("port", "27017")
	this.DbName = section.String("db", "")
	this.ShardBaseNum = section.Int("shard_base_num", this.ShardBaseNum)
	this.User = section.String("user", "")
	this.Pass = section.String("pass", "")
	this.ReplicaSet = section.String("replicaSet", "")
	if this.Host == "" ||
		this.Port == "" ||
		this.Kind == "" ||
		this.DbName == "" {
		panic("required field missing")
	}

	log.Debug("mongodb server: %+v", *this)
}

func (this *ConfigMongodbServer) Address() string {
	addr := "mongodb://" + this.Host + ":" + this.Port + "/"
	if this.ReplicaSet != "" {
		addr += "?replicaSet=" + this.ReplicaSet
	}
	return addr
}

type ConfigMongodb struct {
	ShardBaseNum          int
	Timeout               int
	MaxIdleConnsPerServer int
	HeartbeatInterval     int
	Servers               map[string]*ConfigMongodbServer // key is kind
}

func (this *ConfigMongodb) loadConfig(cf *conf.Conf) {
	this.ShardBaseNum = cf.Int("shard_base_num", 100000)
	this.Timeout = cf.Int("timeout", 30)
	this.MaxIdleConnsPerServer = cf.Int("max_idle_conns_per_server", 2)
	this.HeartbeatInterval = cf.Int("heartbeat_interval", 120)
	this.Servers = make(map[string]*ConfigMongodbServer)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMongodbServer)
		server.ShardBaseNum = this.ShardBaseNum
		server.loadConfig(section)
		this.Servers[server.Kind] = server
	}

	log.Debug("mongodb: %+v", *this)
}
