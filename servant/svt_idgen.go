package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// Ticket server
func (this *FunServantImpl) IdNext(ctx *rpc.Context) (r int64,
	backwards *rpc.TIdTimeBackwards, appErr error) {
	const IDENT = "id.next"

	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	r, appErr = this.idgen.Next()
	if appErr != nil {
		log.Error("Q=%s %s: clock backwards", IDENT, ctx.String())

		backwards = appErr.(*rpc.TIdTimeBackwards)
		appErr = nil
	}

	profiler.do(IDENT, ctx, "{r^%d}", r)

	return
}

func (this *FunServantImpl) IdNextWithTag(ctx *rpc.Context,
	tag int16) (r int64, appErr error) {
	const IDENT = "id.nextag"
	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	r, appErr = this.idgen.NextWithTag(tag)

	profiler.do(IDENT, ctx, "{tag^%d} {r^%d}", tag, r)

	return
}
