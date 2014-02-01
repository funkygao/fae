package memcache

import (
	"github.com/funkygao/fxi/config"
	"hash/crc32"
)

type StandardHash struct {
}

func (this *StandardHash) FindServer(key string) *MemcacheClient {
	servers := config.Servants.Memcache.Servers
	if len(servers) == 1 {
		if memcachePool.BucketCount() == 0 {
			memcachePool.buckets[0] = newMemcacheClient()
			memcachePool.buckets[0].Connect("servers[0].Address()")
		}

		return memcachePool.buckets[0]
	}

	h := crc32.NewIEEE()
	h.Write([]byte(key))
	hash := ((h.Sum32() >> 16) & 0x7fff) % uint32(len(config.Servants.Memcache.Servers))
	if client, present := memcachePool.buckets[hash]; present {
		return client
	}

	// lazy connect to memcached
	memcachePool.buckets[hash] = newMemcacheClient()
	memcachePool.buckets[hash].Connect("addr")
	return memcachePool.buckets[hash]
}

func (this *StandardHash) AddServer(addr string) {

}
