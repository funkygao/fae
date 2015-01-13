package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigGame struct {
	NamegenLength   int
	LockMaxItems    int
	LockExpires     time.Duration
	RedisServerAddr string
	ShardSplit      ConfigGameShardSplit
}

func (this *ConfigGame) LoadConfig(cf *conf.Conf) {
	this.RedisServerAddr = cf.String("redis_server_addr", "127.0.0.1:6379")
	this.NamegenLength = cf.Int("namegen_length", 3)
	this.LockMaxItems = cf.Int("lock_max_items", 1<<20)
	this.LockExpires = cf.Duration("lock_expires", time.Second*10)
	section, err := cf.Section("shard_split_strategy")
	if err != nil {
		panic("empty shard_split_strategy")
	}
	this.ShardSplit.loadConfig(section)

	log.Debug("game conf: %+v", *this)
}

type ConfigGameShardSplit struct {
	Kingdom  int
	User     int
	Alliance int // how many alliances per shard
}

func (this *ConfigGameShardSplit) loadConfig(cf *conf.Conf) {
	this.Kingdom = cf.Int("kingdom", 18000)
	this.User = cf.Int("user", 200000)
	this.Alliance = cf.Int("alliance", 200000/50) // TODO
}
