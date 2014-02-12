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

func (this *FunServantImpl) McSet(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, appErr error) {
	profiler := this.profiler()

	appErr = this.mc.Set(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if appErr == nil {
		r = true
	} else {
		log.Error(appErr)
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

func (this *FunServantImpl) McGet(ctx *rpc.Context, key string) (r []byte,
	miss *rpc.TCacheMissed, appErr error) {
	profiler := this.profiler()

	it, err := this.mc.Get(key)
	if err == nil {
		// cache hit
		r = it.Value
	} else if err == memcache.ErrCacheMiss {
		// cache miss
		miss = rpc.NewTCacheMissed()
		miss.Message = thrift.StringPtr(err.Error()) // optional
	} else {
		appErr = err
		log.Error(err)
	}

	profiler.do("mc.get", ctx,
		"{key^%s} {miss^%v val^%s}",
		key,
		miss,
		r)

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, appErr error) {
	profiler := this.profiler()

	appErr = this.mc.Add(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if appErr == nil {
		r = true
	} else if appErr == memcache.ErrNotStored {
		r = false
		appErr = nil
	} else {
		r = false
		log.Error(appErr)
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
	profiler := this.profiler()

	appErr = this.mc.Delete(key)
	if appErr == nil {
		r = true
	} else if appErr == memcache.ErrCacheMiss {
		r = false
		appErr = nil
	} else {
		log.Error(appErr)
	}

	profiler.do("mc.del", ctx,
		"{key^%s} {err^%s r^%v}",
		key,
		appErr,
		r)

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, key string,
	delta int32) (r int32, appErr error) {
	profiler := this.profiler()

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

	profiler.do("mc.inc", ctx,
		"{key^%s delta^%d} {err^%v r^%d}",
		key,
		delta,
		appErr,
		r)

	return
}
