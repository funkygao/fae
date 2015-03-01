package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/memcache"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/thrift/lib/go/thrift"
)

func (this *FunServantImpl) McSet(ctx *rpc.Context, pool string, key string,
	value *rpc.TMemcacheData, expiration int32) (r bool, ex error) {
	const IDENT = "mc.set"

	if this.mc == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	ex = this.mc.Set(pool, &memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if ex == nil {
		r = true
	} else {
		log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, ex)
	}

	profiler.do(IDENT, ctx,
		"{key^%s val^%s exp^%d} {err^%v r^%v}",
		key,
		value,
		expiration,
		ex,
		r)

	return
}

func (this *FunServantImpl) McGet(ctx *rpc.Context, pool string,
	key string) (r *rpc.TMemcacheData,
	miss *rpc.TCacheMissed, ex error) {
	const IDENT = "mc.get"

	if this.mc == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
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
		ex = err
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
	expiration int32) (r bool, ex error) {
	const IDENT = "mc.add"

	if this.mc == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	ex = this.mc.Add(pool, &memcache.Item{Key: key,
		Value: value.Data, Flags: uint32(value.Flags),
		Expiration: expiration})
	if ex == nil {
		r = true
	} else {
		if ex == memcache.ErrNotStored {
			ex = nil
		} else {
			log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, ex)
		}
	}

	profiler.do(IDENT, ctx,
		"{key^%s val^%s exp^%d} {err^%v r^%v}",
		key,
		value,
		expiration,
		ex,
		r)

	return
}

func (this *FunServantImpl) McDelete(ctx *rpc.Context, pool string,
	key string) (r bool, ex error) {
	const IDENT = "mc.del"

	if this.mc == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	ex = this.mc.Delete(pool, key)
	if ex == nil {
		r = true
	} else {
		if ex == memcache.ErrCacheMiss {
			ex = nil
		} else {
			log.Error("Q=%s %s {key^%s}: %v", IDENT, ctx.String(), key, ex)
		}
	}

	profiler.do(IDENT, ctx,
		"{key^%s} {err^%v r^%v}",
		key,
		ex,
		r)

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, pool string,
	key string, delta int64) (r int64, ex error) {
	const IDENT = "mc.inc"

	if this.mc == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
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
		ex,
		r)

	return
}
