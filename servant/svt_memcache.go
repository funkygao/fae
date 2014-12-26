package servant

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/memcache"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) McSet(ctx *rpc.Context, pool string, key string,
	value *rpc.TMemcacheData, expiration int32) (r bool, appErr error) {
	const IDENT = "mc.set"

	if this.mc == nil {
		appErr = ErrServantNotStarted
		return
	}

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	appErr = this.mc.Set(pool, &memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if appErr == nil {
		r = true
	} else {
		log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx,
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
	const IDENT = "mc.get"

	if this.mc == nil {
		appErr = ErrServantNotStarted
		return
	}

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

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
		log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, err)
	}

	profiler.do(IDENT, ctx,
		"{key^%s} {miss^%v val^%s}",
		key,
		miss,
		r)

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.Context, pool string, key string,
	value *rpc.TMemcacheData,
	expiration int32) (r bool, appErr error) {
	const IDENT = "mc.add"

	if this.mc == nil {
		appErr = ErrServantNotStarted
		return
	}

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	appErr = this.mc.Add(pool, &memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if appErr == nil {
		r = true
	} else {
		if appErr == memcache.ErrNotStored {
			appErr = nil
		} else {
			log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, appErr)
		}
	}

	profiler.do(IDENT, ctx,
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
	const IDENT = "mc.del"

	if this.mc == nil {
		appErr = ErrServantNotStarted
		return
	}

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	appErr = this.mc.Delete(pool, key)
	if appErr == nil {
		r = true
	} else {
		if appErr == memcache.ErrCacheMiss {
			appErr = nil
		} else {
			log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, appErr)
		}
	}

	profiler.do(IDENT, ctx,
		"{key^%s} {err^%v r^%v}",
		key,
		appErr,
		r)

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, pool string,
	key string, delta int64) (r int64, appErr error) {
	const IDENT = "mc.inc"

	if this.mc == nil {
		appErr = ErrServantNotStarted
		return
	}

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	newVal, err := this.mc.Increment(pool, key, delta)
	if err == nil {
		r = int64(newVal)
	} else if err != memcache.ErrCacheMiss {
		log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, err)
	}

	profiler.do(IDENT, ctx,
		"{key^%s delta^%d} {err^%v r^%d}",
		key,
		delta,
		appErr,
		r)

	return
}
