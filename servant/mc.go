package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fxi/servant/memcache"
	"time"
)

func (this *FunServantImpl) McSet(key string, value []byte, expiration int32) (r bool, err error) {
	this.t1 = time.Now()
	err = this.mc.Set(&memcache.Item{Key: key, Value: value, Expiration: expiration})
	if err == nil {
		r = true
	}

	log.Debug("mc_set key:%s value:%s, expiration:%v %s", key, string(value), expiration,
		time.Since(this.t1))

	return
}

func (this *FunServantImpl) McGet(key string) (r []byte, err error) {
	this.t1 = time.Now()
	var it *memcache.Item
	it, err = this.mc.Get(key)
	if err == nil {
		r = it.Value
	}

	log.Debug("mc_get key:%s %s", key, time.Since(this.t1))
	return
}
