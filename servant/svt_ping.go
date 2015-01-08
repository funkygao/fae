package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/server"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, appErr error) {
	const IDENT = "ping"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	r = fmt.Sprintf("pong, %s, myid:%d", server.BuildID, this.conf.IdgenWorkerId)

	profiler.do(IDENT, ctx, "pong")
	return
}
