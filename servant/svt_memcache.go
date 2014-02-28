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

func (this *FunServantImpl) McSet(ctx *rpc.Context, key string,
	value *rpc.TMemcacheData, expiration int32) (r bool, appErr error) {
	this.stats.inc("mc.set")
	this.stats.inBytes.Inc(int64(len(value.Data) + len(key)))

	profiler := this.profiler()
	err := this.mc.Set(&memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if err == nil {
		r = true
	} else {
		log.Error("mc.set: %v", err)
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

func (this *FunServantImpl) McGet(ctx *rpc.Context,
	key string) (r *rpc.TMemcacheData,
	miss *rpc.TCacheMissed, appErr error) {
	this.stats.inc("mc.get")
	this.stats.inBytes.Inc(int64(len(key)))

	profiler := this.profiler()
	it, err := this.mc.Get(key)
	if err == nil {
		// cache hit
		r = rpc.NewTMemcacheData()
		r.Data = it.Value
		r.Flags = int32(it.Flags)
		this.stats.outBytes.Inc(int64(len(r.Data)))
	} else if err == memcache.ErrCacheMiss {
		// cache miss
		miss = rpc.NewTCacheMissed()
		miss.Message = thrift.StringPtr(err.Error()) // optional
	} else {
		log.Error("mc.get: %v", err)
	}

	profiler.do("mc.get", ctx,
		"{key^%s} {miss^%v val^%s}",
		key,
		miss,
		r)

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.Context, key string,
	value *rpc.TMemcacheData,
	expiration int32) (r bool, appErr error) {
	this.stats.inc("mc.add")
	this.stats.inBytes.Inc(int64(len(key) + len(value.Data)))

	profiler := this.profiler()
	err := this.mc.Add(&memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if err == nil {
		r = true
	} else {
		log.Error("mc.add[%s]: %v", key, err)
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

func (this *FunServantImpl) McDelete(ctx *rpc.Context, key string) (r bool,
	appErr error) {
	this.stats.inc("mc.del")
	this.stats.inBytes.Inc(int64(len(key)))

	profiler := this.profiler()
	err := this.mc.Delete(key)
	if err == nil {
		r = true
	} else {
		log.Error("mc.del: %v", err)
	}

	profiler.do("mc.del", ctx,
		"{key^%s} {err^%v r^%v}",
		key,
		appErr,
		r)

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, key string,
	delta int32) (r int32, appErr error) {
	this.stats.inc("mc.inc")
	this.stats.inBytes.Inc(int64(len(key)))

	var (
		newVal uint64
		err    error
	)
	profiler := this.profiler()
	if delta > 0 {
		newVal, err = this.mc.Increment(key, uint64(delta))
	} else {
		newVal, err = this.mc.Decrement(key, uint64(delta))
	}

	if err == nil {
		r = int32(newVal)
	}

	profiler.do("mc.inc", ctx,
		"{key^%s delta^%d} {err^%v r^%d}",
		key,
		delta,
		appErr,
		r)

	return
}
