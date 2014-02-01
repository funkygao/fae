package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fxi/config"
	"github.com/funkygao/fxi/servant/memcache"
	"github.com/funkygao/golib/cache"
	"time"
)

type FunServantImpl struct {
	conf *config.ConfigServant

	t1 time.Time // timeit

	mc *memcache.Client
	lc *cache.LruCache
}

func NewFunServant(cf *config.ConfigServant) (this *FunServantImpl) {
	this = &FunServantImpl{conf: cf}
	this.lc = cache.NewLruCache(1 << 30)
	memcacheServers := this.conf.Memcache.ServerList()
	this.mc = memcache.New(this.conf.Memcache.HashStrategy, memcacheServers...)

	log.Debug("memcache servers %v", memcacheServers)
	return
}

func (this *FunServantImpl) Ping() (r string, err error) {
	log.Debug("ping")
	return "pong", nil
}
