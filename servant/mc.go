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

func (this *FunServantImpl) McAdd(ctx *rpc.ReqCtx, key string, value []byte,
	expiration int32) (r bool, err error) {
	e := this.mc.Add(&memcache.Item{Key: key, Value: value, Expiration: expiration})
	if e == nil {
		r = true
	}

	return
}

func (this *FunServantImpl) McDelete(ctx *rpc.ReqCtx, key string) (r bool, err error) {
	e := this.mc.Delete(key)
	if e == nil {
		r = true
	}

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.ReqCtx, key string, delta int32) (r int32, err error) {
	var newVal uint64
	newVal, err = this.mc.Increment(key, uint64(delta))
	if err == nil {
		r = int32(newVal)
	}

	return
}
