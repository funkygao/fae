package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/memcache"
	"time"
)

func (this *FunServantImpl) McSet(ctx *rpc.ReqCtx, key string, value []byte,
	expiration int32) (r bool, err error) {
	this.t1 = time.Now()
	err = this.mc.Set(&memcache.Item{Key: key, Value: value, Expiration: expiration})
	if err == nil {
		r = true
	}

	log.Debug("ctx:%+v mc_set key:%s value:%s, expiration:%v %s",
		*ctx,
		key, string(value), expiration,
		time.Since(this.t1))

	return
}

func (this *FunServantImpl) McGet(ctx *rpc.ReqCtx, key string) (r []byte, err error) {
	this.t1 = time.Now()
	var it *memcache.Item
	it, err = this.mc.Get(key)
	if err == nil {
		r = it.Value
	}

	log.Debug("ctx:%+v mc_get key:%s %s",
		*ctx,
		key, time.Since(this.t1))
	return
}
