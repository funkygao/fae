package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) KvdbSet(ctx *rpc.Context,
	key []byte, value []byte) (r bool, appErr error) {
	this.stats.inc("kvdb.set")

	profiler := this.profiler()
	if err := this.kvdb.Put(key, value); err != nil {
		log.Error("kvdb.set: %v", err)
	} else {
		r = true
	}
	profiler.do("kv.set", ctx,
		"{key^%s val^%s} {r^%v}",
		key, value, r)

	return
}

func (this *FunServantImpl) KvdbGet(ctx *rpc.Context, key []byte) (r []byte,
	appErr error) {
	this.stats.inc("kvdb.get")

	profiler := this.profiler()
	r, _ = this.kvdb.Get(key)
	profiler.do("kv.get", ctx,
		"{key^%s} {r^%s}",
		key, r)
	return
}

func (this *FunServantImpl) KvdbDel(ctx *rpc.Context, key []byte) (r bool,
	appErr error) {
	this.stats.inc("kvdb.del")

	profiler := this.profiler()
	err := this.kvdb.Delete(key)
	if err == nil {
		r = true
	}
	profiler.do("kv.del", ctx,
		"{key^%s} {r^%v}",
		key, r)
	return
}
