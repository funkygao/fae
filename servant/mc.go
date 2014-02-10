/*
MCache key:string, value:[]byte.
*/
package servant

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/memcache"
	"github.com/funkygao/golib/sampling"
	log "github.com/funkygao/log4go"
	"time"
)

func (this *FunServantImpl) McSet(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, intError error) {
	profiler := this.profilerInfo(ctx)

	intError = this.mc.Set(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	if profiler.do {
		log.Debug("T=%s Q=mc.set %s {key^%s val^%s exp^%d} {%v}",
			time.Since(profiler.t1),
			this.callerInfo(ctx),
			key, this.truncValue(value), expiration,
			r)
	}

	return
}

func (this *FunServantImpl) McGet(ctx *rpc.Context, key string) (r []byte,
	miss *rpc.TCacheMissed, intError error) {
	profiler := this.profilerInfo(ctx)

	it, err := this.mc.Get(key)
	if err == nil {
		// cache hit
		r = it.Value
	} else if err == memcache.ErrCacheMiss {
		// cache miss
		if sampling.SampleRateSatisfied(5) {
			log.Debug("mc missed: %s", key)
		}

		miss = rpc.NewTCacheMissed()
		miss.Message = thrift.StringPtr(err.Error()) // optional
	} else {
		intError = err
		log.Error(err)
	}

	if profiler.do {
		log.Debug("T=%s Q=mc.get %s {key^%s} {%s}",
			time.Since(profiler.t1),
			this.callerInfo(ctx),
			key,
			this.truncValue(r))
	}

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, intError error) {
	profiler := this.profilerInfo(ctx)

	intError = this.mc.Add(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if intError == nil {
		r = true
	} else if intError == memcache.ErrNotStored {
		r = false
		intError = nil
	} else {
		r = false
		log.Error(intError)
	}

	if profiler.do {
		log.Debug("T=%s Q=mc.add %s {key^%s val^%s exp^%d} {%v}",
			time.Since(profiler.t1),
			this.callerInfo(ctx),
			key, this.truncValue(value), expiration,
			r)
	}

	return
}

func (this *FunServantImpl) McDelete(ctx *rpc.Context, key string) (r bool,
	intError error) {
	profiler := this.profilerInfo(ctx)

	intError = this.mc.Delete(key)
	if intError == nil {
		r = true
	} else if intError == memcache.ErrCacheMiss {
		r = false
		intError = nil
	} else {
		log.Error(intError)
	}

	if profiler.do {
		log.Debug("T=%s Q=mc.del %s {key^%s} {%v}",
			time.Since(profiler.t1),
			this.callerInfo(ctx),
			key,
			r)
	}

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, key string,
	delta int32) (r int32, intError error) {
	profiler := this.profilerInfo(ctx)

	var (
		newVal uint64
		err    error
	)
	if delta > 0 {
		newVal, err = this.mc.Increment(key, uint64(delta))
	} else {
		newVal, err = this.mc.Decrement(key, uint64(delta))
	}

	if err == nil {
		r = int32(newVal)
	}

	if profiler.do {
		log.Debug("T=%s Q=mc.inc %s {key^%s delta^%d} {%d}",
			time.Since(profiler.t1),
			this.callerInfo(ctx),
			key, delta,
			r)
	}

	return
}
