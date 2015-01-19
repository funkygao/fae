package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/idgen"
	log "github.com/funkygao/log4go"
)

// Ticket server
func (this *FunServantImpl) IdNext(ctx *rpc.Context) (r int64,
	backwards *rpc.TIdTimeBackwards, ex error) {
	const IDENT = "id.next"

	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	r, ex = this.idgen.Next()
	if ex != nil {
		log.Error("Q=%s %s: clock backwards", IDENT, ctx.String())

		backwards = ex.(*rpc.TIdTimeBackwards)
		ex = nil
	}

	profiler.do(IDENT, ctx, "{r^%d}", r)

	return
}

func (this *FunServantImpl) IdNextWithTag(ctx *rpc.Context,
	tag int16) (r int64, ex error) {
	const IDENT = "id.nextag"
	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	r, ex = this.idgen.NextWithTag(tag)

	profiler.do(IDENT, ctx, "{tag^%d} {r^%d}", tag, r)

	return
}

func (this *FunServantImpl) IdDecode(ctx *rpc.Context,
	id int64) (r []int64, ex error) {
	const IDENT = "id.decode"
	this.stats.inc(IDENT)
	ts, tag, wid, seq := idgen.DecodeId(id)
	r = []int64{ts, tag, wid, seq}
	return
}
