/*
MCache key:string, value:[]byte.
*/
package servant

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/memcache"
)

func (this *FunServantImpl) McSet(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, intError error) {
	profiler := this.profiler()

	intError = this.mc.Set(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	profiler.do("mc.set", ctx, "{key^%s val^%s exp^%d} {%v}",
		key,
		this.truncatedBytes(value),
		expiration,
		r)

	return
}

func (this *FunServantImpl) McGet(ctx *rpc.Context, key string) (r []byte,
	miss *rpc.TCacheMissed, intError error) {
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
		intError = err
		log.Error(err)
	}

	profiler.do("mc.get", ctx, "{key^%s} {%s}",
		key,
		this.truncatedBytes(r))

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, intError error) {
	profiler := this.profiler()

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

	profiler.do("mc.add", ctx, "{key^%s val^%s exp^%d} {%v}",
		key,
		this.truncatedBytes(value),
		expiration,
		r)

	return
}

func (this *FunServantImpl) McDelete(ctx *rpc.Context, key string) (r bool,
	intError error) {
	profiler := this.profiler()

	intError = this.mc.Delete(key)
	if intError == nil {
		r = true
	} else if intError == memcache.ErrCacheMiss {
		r = false
		intError = nil
	} else {
		log.Error(intError)
	}

	profiler.do("mc.del", ctx, "{key^%s} {%v}",
		key,
		r)

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, key string,
	delta int32) (r int32, intError error) {
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

	profiler.do("mc.inc", ctx, "{key^%s delta^%d} {%d}",
		key,
		delta,
		r)

	return
}
