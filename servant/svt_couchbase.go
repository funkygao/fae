package servant

import (
	//couchbase "github.com/couchbaselabs/go-couchbase"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// curl localhost:8091/pools/ | python -m json.tool

func (this *FunServantImpl) CbAdd(ctx *rpc.Context, key string, val []byte, expire int) (appErr error) {
	log.Debug("cb")

	return nil
}

func (this *FunServantImpl) CbSet(ctx *rpc.Context, key string) (appErr error) {

	return nil
}

func (this *FunServantImpl) CbDel(ctx *rpc.Context, key string) (appErr error) {

	return nil
}

func (this *FunServantImpl) CbGet(ctx *rpc.Context, key string) (appErr error) {

	return nil
}

func (this *FunServantImpl) CbGetBulk(ctx *rpc.Context, key string) (appErr error) {

	return nil
}

func (this *FunServantImpl) CbInc(ctx *rpc.Context, key string) (appErr error) {

	return nil
}
