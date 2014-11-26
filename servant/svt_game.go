package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/gofmt"
	log "github.com/funkygao/log4go"
)

// get a uniq name with length 3
func (this *FunServantImpl) GmName3(ctx *rpc.Context) (r string, appErr error) {
	const IDENT = "gm.name3"

	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	r = this.namegen.Next()

	profiler.do(IDENT, ctx, "{r^%s}", r)

	return
}

// record php request time and payload size in bytes
func (this *FunServantImpl) GmLatency(ctx *rpc.Context, ms int32,
	bytes int32) (appErr error) {
	this.phpLatency.Update(int64(ms))
	this.phpPayloadSize.Update(int64(bytes))

	log.Trace("{%dms %s}: {uid^%d rid^%s reason^%s}: ",
		ms, gofmt.ByteSize(bytes),
		ctx.Uid, ctx.Rid, ctx.Reason)

	return
}
