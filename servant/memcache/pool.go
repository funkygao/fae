package memcache

import (
	"github.com/funkygao/fxi/config"
	"sync"
)

type MemcachePool struct {
	*sync.RWMutex

	hash Hash

	buckets map[uint32]*MemcacheClient // hash -> client
}

func newMemcachePool() (this *MemcachePool) {
	this = new(MemcachePool)
	this.RWMutex = new(sync.RWMutex)
	return
}

func (this *MemcachePool) BucketCount() int {
	return len(this.buckets)
}

func (this *MemcachePool) Init(cf *config.ConfigMemcache) {
	switch cf.HashStrategy {
	case "standard":
		this.hash = new(StandardHash)

	case "consistent":

	default:
		panic("Invalid hash_strategy: " + cf.HashStrategy)
	}

	this.buckets = make(map[uint32]*MemcacheClient)
}

func (this *MemcachePool) AddServer(addr string) {

}
