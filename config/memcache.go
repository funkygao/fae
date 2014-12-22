package config

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigMemcacheServer struct {
	Pool string
	Host string
	Port string
}

func (this *ConfigMemcacheServer) loadConfig(section *conf.Conf) {
	this.Host = section.String("host", "")
	if this.Host == "" {
		panic("Empty memcache server host")
	}
	this.Port = section.String("port", "")
	if this.Port == "" {
		panic("Empty memcache server port")
	}
	this.Pool = section.String("pool", "default")

}

func (this *ConfigMemcacheServer) Address() string {
	return this.Host + ":" + this.Port
}

type ConfigMemcache struct {
	HashStrategy string
	// for both conn and io timeout
	Timeout               time.Duration
	MaxIdleConnsPerServer int
	MaxConnsPerServer     int
	ReplicaN              int
	Breaker               ConfigBreaker
	Servers               map[string]*ConfigMemcacheServer // key is host:port(addr)
}

func (this *ConfigMemcache) ServerList() []string {
	servers := make([]string, len(this.Servers))
	i := 0
	for addr, _ := range this.Servers {
		servers[i] = addr
		i += 1
	}

	return servers
}

func (this *ConfigMemcache) Pools() (pools []string) {
	poolsMap := make(map[string]bool)
	for _, server := range this.Servers {
		poolsMap[server.Pool] = true
	}
	for poolName, _ := range poolsMap {
		pools = append(pools, poolName)
	}
	return
}

func (this *ConfigMemcache) Enabled() bool {
	return len(this.Servers) > 0
}

func (this *ConfigMemcache) loadConfig(cf *conf.Conf) {
	this.Servers = make(map[string]*ConfigMemcacheServer)
	this.HashStrategy = cf.String("hash_strategy", "standard")
	this.Timeout = cf.Duration("timeout", 4*time.Second)
	this.ReplicaN = cf.Int("replica_num", 1)
	section, err := cf.Section("breaker")
	if err == nil {
		this.Breaker.loadConfig(section)
	}
	this.MaxIdleConnsPerServer = cf.Int("max_idle_conns_per_server", 3)
	this.MaxConnsPerServer = cf.Int("max_conns_per_server",
		this.MaxIdleConnsPerServer*10)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMemcacheServer)
		server.loadConfig(section)
		this.Servers[server.Address()] = server
	}

	log.Debug("memcache conf: %+v", *this)
}
