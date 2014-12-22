package config

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigRedisServer struct {
	Addr        string // host:port
	MaxIdle     int
	IdleTimeout time.Duration
}

func (this *ConfigRedisServer) loadConfig(cf *conf.Conf) {
	this.Addr = cf.String("addr", "")
	if this.Addr == "" {
		panic("Empty redis server addr")
	}
	this.MaxIdle = cf.Int("max_idle", 10)
	this.IdleTimeout = cf.Duration("idle_timeout", 10*time.Minute)
}

type ConfigRedis struct {
	Breaker ConfigBreaker
	Servers map[string]map[string]*ConfigRedisServer // pool:serverAddr:ConfigRedisServer
}

func (this *ConfigRedis) loadConfig(cf *conf.Conf) {
	section, err := cf.Section("breaker")
	if err == nil {
		this.Breaker.loadConfig(section)
	}

	for i := 0; i < len(cf.List("pools", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("pools[%d]", i))
		if err != nil {
			panic(err)
		}

		pool := section.String("name", "")
		if pool == "" {
			panic("Empty redis pool name")
		}

		this.Servers[pool] = make(map[string]*ConfigRedisServer)

		// get servers in each pool
		for j := 0; j < len(section.List("servers", nil)); j++ {
			server, err := section.Section(fmt.Sprintf("servers[%d]", j))
			if err != nil {
				panic(err)
			}

			redisServer := new(ConfigRedisServer)
			redisServer.loadConfig(server)
			this.Servers[pool][redisServer.Addr] = redisServer
		}
	}

	log.Debug("redis conf: %+v", *this)
}

func (this *ConfigRedis) PoolServers(pool string) []string {
	r := make([]string, 0)
	for addr, _ := range this.Servers[pool] {
		r = append(r, addr)
	}
	return r
}

func (this *ConfigRedis) Enabled() bool {
	return len(this.Servers) > 0
}
