package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, appErr error) {
	const IDENT = "ping"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	profiler.do(IDENT, ctx, "pong")
	return "pong", nil
}
