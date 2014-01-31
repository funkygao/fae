package config

import (
	log "code.google.com/p/log4go"
	"fmt"
	conf "github.com/funkygao/jsconf"
)

type ConfigMongodbServer struct {
	ShardName  string
	Host       string
	Port       int
	User       string
	Pass       string
	DbName     string
	ReplicaSet string
}

func (this *ConfigMongodbServer) loadConfig(section *conf.Conf) {
	this.ShardName = section.String("shard_name", "")
	this.Host = section.String("host", "")
	this.Port = section.Int("port", 0)
	this.DbName = section.String("db", "")
	this.User = section.String("user", "")
	this.Pass = section.String("pass", "")
	this.ReplicaSet = section.String("replicaSet", "")
	if this.Host == "" ||
		this.Port == 0 ||
		this.ShardName == "" ||
		this.DbName == "" {
		panic("required filed")
	}

	log.Debug("mongodb server: %+v", *this)
}

type ConfigMongodb struct {
	Servers map[string]*ConfigMongodbServer // key is shardName
}

func (this *ConfigMongodb) loadConfig(cf *conf.Conf) {
	this.Servers = make(map[string]*ConfigMongodbServer)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMongodbServer)
		server.loadConfig(section)
		this.Servers[server.ShardName] = server
	}

	log.Debug("mongodb: %+v", *this)
}
