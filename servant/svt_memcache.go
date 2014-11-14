/*
memcache key:string, value:[]byte.
*/
package servant

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/memcache"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) McSet(ctx *rpc.Context, pool string, key string,
	value *rpc.TMemcacheData, expiration int32) (r bool, appErr error) {
	this.stats.inc("mc.set")

	profiler := this.getSession(ctx).startProfiler()
	appErr = this.mc.Set(pool, &memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if appErr == nil {
		r = true
	} else {
		log.Error("mc.set {key^%s}: %v", key, appErr)
	}

	profiler.do("mc.set", ctx,
		"{key^%s val^%s exp^%d} {err^%v r^%v}",
		key,
		value,
		expiration,
		appErr,
		r)

	return
}

func (this *FunServantImpl) McGet(ctx *rpc.Context, pool string,
	key string) (r *rpc.TMemcacheData,
	miss *rpc.TCacheMissed, appErr error) {
	this.stats.inc("mc.get")

	profiler := this.getSession(ctx).startProfiler()
	it, err := this.mc.Get(pool, key)
	if err == nil {
		// cache hit
		r = rpc.NewTMemcacheData()
		r.Data = it.Value
		r.Flags = int32(it.Flags)
	} else if err == memcache.ErrCacheMiss {
		// cache miss
		miss = rpc.NewTCacheMissed()
		miss.Message = thrift.StringPtr(err.Error()) // optional
	} else {
		appErr = err
		log.Error("mc.get {key^%s}: %v", key, err)
	}

	profiler.do("mc.get", ctx,
		"{key^%s} {miss^%v val^%s}",
		key,
		miss,
		r)

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.Context, pool string, key string,
	value *rpc.TMemcacheData,
	expiration int32) (r bool, appErr error) {
	this.stats.inc("mc.add")

	profiler := this.getSession(ctx).startProfiler()
	appErr = this.mc.Add(pool, &memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if appErr == nil {
		r = true
	} else {
		if appErr == memcache.ErrNotStored {
			appErr = nil
		} else {
			log.Error("mc.add {key^%s}: %v", key, appErr)
		}
	}

	profiler.do("mc.add", ctx,
		"{key^%s val^%s exp^%d} {err^%v r^%v}",
		key,
		value,
		expiration,
		appErr,
		r)

	return
}

func (this *FunServantImpl) McDelete(ctx *rpc.Context, pool string,
	key string) (r bool, appErr error) {
	this.stats.inc("mc.del")

	profiler := this.getSession(ctx).startProfiler()
	appErr = this.mc.Delete(pool, key)
	if appErr == nil {
		r = true
	} else {
		if appErr == memcache.ErrCacheMiss {
			appErr = nil
		} else {
			log.Error("mc.del {key^%s}: %v", key, appErr)
		}
	}

	profiler.do("mc.del", ctx,
		"{key^%s} {err^%v r^%v}",
		key,
		appErr,
		r)

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, pool string,
	key string, delta int64) (r int64, appErr error) {
	this.stats.inc("mc.inc")

	profiler := this.getSession(ctx).startProfiler()

	newVal, err := this.mc.Increment(pool, key, delta)
	if err == nil {
		r = int64(newVal)
	} else if err != memcache.ErrCacheMiss {
		log.Error("mc.inc {key^%s}: %v", key, err)
	}

	profiler.do("mc.inc", ctx,
		"{key^%s delta^%d} {err^%v r^%d}",
		key,
		delta,
		appErr,
		r)

	return
}
