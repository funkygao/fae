package memcache

import (
//	"github.com/funkygao/fxi/config"
)

type FindServer func(key string) *MemcacheClient

func standardFindServer(key string) *MemcacheClient {
	return nil
}

func consistentFindServer(key string) *MemcacheClient {
	return nil
}
