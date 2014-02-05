/*
MCache key:string, value:[]byte.
*/
package servant

import (
	log "code.google.com/p/log4go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/memcache"
)

func (this *FunServantImpl) McSet(ctx *rpc.ReqCtx, key string, value []byte,
	expiration int32) (r bool, intError error) {
	intError = this.mc.Set(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	return
}

func (this *FunServantImpl) McGet(ctx *rpc.ReqCtx, key string) (r []byte,
	miss *rpc.TCacheMissed, intError error) {
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

	return
}

func (this *FunServantImpl) McAdd(ctx *rpc.ReqCtx, key string, value []byte,
	expiration int32) (r bool, intError error) {
	intError = this.mc.Add(&memcache.Item{Key: key, Value: value,
		Expiration: expiration})
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	return
}

func (this *FunServantImpl) McDelete(ctx *rpc.ReqCtx, key string) (r bool,
	intError error) {
	intError = this.mc.Delete(key)
	if intError == nil {
		r = true
	} else {
		log.Error(intError)
	}

	return
}

func (this *FunServantImpl) McIncrement(ctx *rpc.ReqCtx, key string,
	delta int32) (r int32, intError error) {
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
