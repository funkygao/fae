package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// curl localhost:8091/pools/ | python -m json.tool

func (this *FunServantImpl) CbGet(ctx *rpc.Context, bucket string,
	key string) (r []byte, appErr error) {
	const IDENT = "cb.get"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	pool, _ := this.cb.GetPool("default")
	b, _ := pool.GetBucket(bucket)

	r, appErr = b.GetRaw(key)
	if appErr != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx,
		"{bucket^%s key^%s} {r^%s}",
		bucket, key, string(r))

	return
}

func (this *FunServantImpl) CbSet(ctx *rpc.Context, bucket string,
	key string, val []byte, expire int32) (r bool, appErr error) {
	const IDENT = "cb.set"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	pool, _ := this.cb.GetPool("default")
	b, _ := pool.GetBucket(bucket)

	appErr = b.SetRaw(key, int(expire), val)
	if appErr == nil {
		r = true
	}

	profiler.do(IDENT, ctx,
		"{bucket^%s key^%s} {r^%v}",
		bucket, key, r)

	return
}
