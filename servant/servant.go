package servant

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/memcache"
	"github.com/funkygao/fae/servant/mongo"
	"github.com/funkygao/golib/cache"
	"time"
)

type FunServantImpl struct {
	conf *config.ConfigServant

	lc *cache.LruCache
	mc *memcache.Client
	mg *mongo.Client
}

func NewFunServant(cf *config.ConfigServant) (this *FunServantImpl) {
	this = &FunServantImpl{conf: cf}
	this.lc = cache.NewLruCache(this.conf.Lcache.LruMaxItems)
	this.lc.OnEvicted = func(key cache.Key, value interface{}) {
	}

	memcacheServers := this.conf.Memcache.ServerList()
	this.mc = memcache.New(this.conf.Memcache.HashStrategy, memcacheServers...)
	this.mc.Timeout = time.Duration(this.conf.Memcache.Timeout) * time.Second
	this.mc.MaxIdleConnsPerServer = this.conf.Memcache.MaxIdleConnsPerServer

	this.mg = mongo.New(this.conf.Mongodb)

	go this.runWatchdog()

	return
}
