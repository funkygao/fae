package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) RdCall(ctx *rpc.Context, cmd string,
	pool string, key string, args []string) (r string, appErr error) {
	const IDENT = "rd.call"

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	var val interface{}
	if val, appErr = this.rd.Call(cmd, pool, key, args...); appErr == nil {
		r = val.(string)
	}

	profiler.do(IDENT, ctx,
		"{cmd^%s pool^%s key^%s arg^%+v} {err^%v r^%s}",
		cmd, pool, key, args,
		appErr, r)

	return
}
