package memcache

import (
	"github.com/funkygao/fxi/config"
	"sync"
)

type MemcachePool struct {
	*sync.RWMutex

	findServer FindServer

	servers map[uint32]*MemcacheClient // hash -> client
}

func newMemcachePool() (this *MemcachePool) {
	this = new(MemcachePool)
	this.RWMutex = new(sync.RWMutex)
	return
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

	this.servers = make(map[uint32]*MemcacheClient)
}
