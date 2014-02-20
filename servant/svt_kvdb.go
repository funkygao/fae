package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) KvdbSet(ctx *rpc.Context,
	key []byte, value []byte) (r bool, appErr error) {
	this.stats.inc("kvdb.set")
	if err := this.kvdb.Put(key, value); err != nil {
		log.Error("kvdb.set: %v", err)
	} else {
		r = true
	}

	return
}

func (this *FunServantImpl) KvdbGet(ctx *rpc.Context, key []byte) (r []byte,
	appErr error) {
	this.stats.inc("kvdb.get")
	r, _ = this.kvdb.Get(key)
	return
}

func (this *FunServantImpl) KvdbDel(ctx *rpc.Context, key []byte) (appErr error) {
	this.stats.inc("kvdb.del")
	this.kvdb.Delete(key)
	return
}
