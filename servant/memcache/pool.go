package memcache

import (
	"github.com/funkygao/fxi/config"
)

type MemcachePool struct {
	findServer FindServer

	servers map[string]*MemcacheClient
}

func (this *MemcachePool) Init(cf *config.ConfigMemcache) {
	switch cf.HashStrategy {
	case "standard":
		this.findServer = standardFindServer

	case "consistent":
		this.findServer = consistentFindServer

	default:
		panic("Invalid hash_strategy: " + cf.HashStrategy)
	}

	this.servers = make(map[string]*MemcacheClient)
	for addr, _ := range cf.Servers {
		this.servers[addr] = newMemcacheClient()
		this.servers[addr].Connect(addr)
	}

}
