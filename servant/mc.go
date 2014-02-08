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
)

func (this *FunServantImpl) McSet(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, intError error) {
	log.Debug("mc_set %s> key=%s value=%s", this.callerInfo(ctx),
		key, value)

	intError = this.mc.Set(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	return
}

func (this *FunServantImpl) McGet(ctx *rpc.Context, key string) (r []byte,
	miss *rpc.TCacheMissed, intError error) {
	log.Debug("mc_get %s> key=%s", this.callerInfo(ctx), key)

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

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.Context, key string, value []byte,
	expiration int32) (r bool, intError error) {
	log.Debug("mc_add %s> key=%s value=%s", this.callerInfo(ctx),
		key, value)

	intError = this.mc.Add(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	return
}

func (this *FunServantImpl) McDelete(ctx *rpc.Context, key string) (r bool,
	intError error) {
	intError = this.mc.Delete(key)
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.Context, key string,
	delta int32) (r int32, intError error) {
	log.Debug("mc_inc %s> key=%s delta=%d", this.callerInfo(ctx),
		key, delta)

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

	return
}
