package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// curl localhost:8091/pools/ | python -m json.tool
// curl localhost:8091/poolsStreaming/default?uuid=ee6009fb8f1ba20b3101a465455828ee

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

	b, _ := this.cb.GetBucket(bucket)

	r, appErr = b.GetRaw(key) // FIXME 如果不存在，也会抛错，需要额外处理
	if appErr != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s} {r^%s}",
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

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		log.Error(err)
	}

	appErr = b.SetRaw(key, int(expire), val)
	if appErr == nil {
		r = true
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s v^%s} {r^%v}",
		bucket, key, string(val), r)

	return
}
