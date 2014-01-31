package memcache

import (
	"github.com/funkygao/fxi/config"
	"hash/crc32"
)

type FindServer func(key string) *MemcacheClient

func standardFindServer(key string) *MemcacheClient {
	if len(config.Servants.Memcache.Servers) == 1 {

	}

	h := crc32.NewIEEE()
	h.Write([]byte(key))
	hash := (h.Sum32() >> 16) & 0x7fff
	bucket := hash % uint32(len(config.Servants.Memcache.Servers))
	if client, present := memcachePool.servers[bucket]; present {
		return client
	}

	// lazy connect to memcached
	memcachePool.servers[bucket] = newMemcacheClient()
	memcachePool.servers[bucket].Connect("addr")
}

func consistentFindServer(key string) *MemcacheClient {
	return nil
}
