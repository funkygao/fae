/*
LCache key:string, value:[]byte.
*/
package servant

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) LcSet(ctx *rpc.ReqCtx,
	key string, value []byte) (r bool, intError error) {
	this.lc.Set(key, value)
	r = true

	return
}

func (this *FunServantImpl) LcGet(ctx *rpc.ReqCtx, key string) (r []byte,
	miss *rpc.TCacheMissed, intError error) {
	result, ok := this.lc.Get(key)
	if !ok {
		miss = rpc.NewTCacheMissed()
		miss.Message = thrift.StringPtr("lcache missed: " + key) // optional
	} else {
		r = result.([]byte)
	}

	return
}

func (this *FunServantImpl) LcDel(ctx *rpc.ReqCtx, key string) (intError error) {
	this.lc.Del(key)
	return
}
