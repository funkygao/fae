package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/gofmt"
	log "github.com/funkygao/log4go"
)

// get a uniq name with length 3
// TODO dump to redis periodically
func (this *FunServantImpl) GmName3(ctx *rpc.Context) (r string, appErr error) {
	const IDENT = "gm.name3"

	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	r = this.namegen.Next()

	// replication of name to peers in cluster in async mode
	go func() {
		for _, svt := range this.proxy.ClusterServants() {
			log.Debug("%s: %s -> %s", IDENT, r, svt.Addr())

			svt.HijackContext(ctx)
			svt.SyncName3(ctx, r)
			svt.Recycle() // VERY important
		}
	}()

	profiler.do(IDENT, ctx, "{r^%s}", r)

	return
}

// record php request time and payload size in bytes
func (this *FunServantImpl) GmLatency(ctx *rpc.Context, ms int32,
	bytes int32) (appErr error) {
	this.phpLatency.Update(int64(ms))
	this.phpPayloadSize.Update(int64(bytes))

	log.Trace("{%dms %s}: {uid^%d rid^%s reason^%s}",
		ms, gofmt.ByteSize(bytes),
		this.extractUid(ctx), ctx.Rid, ctx.Reason)

	return
}
