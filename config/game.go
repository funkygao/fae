package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigGame struct {
	NamegenLength        int
	LockMaxItems         int
	LockExpires          time.Duration
	MaxPlayersPerKingdom int
	RedisServerAddr      string
}

func (this *ConfigGame) LoadConfig(cf *conf.Conf) {
	this.RedisServerAddr = cf.String("redis_server_addr", "127.0.0.1:6379")
	this.MaxPlayersPerKingdom = cf.Int("max_players_per_kingdom", 18000)
	this.NamegenLength = cf.Int("namegen_length", 3)
	this.LockMaxItems = cf.Int("lock_max_items", 1<<20)
	this.LockExpires = cf.Duration("lock_expires", time.Second*10)
	log.Debug("game conf: %+v", *this)
}
