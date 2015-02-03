package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/server"
	"time"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, ex error) {
	const IDENT = "ping"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	svtStats.inc(IDENT)

	r = fmt.Sprintf("ver:%s, build:%s, myid:%d, uptime:%s",
		server.VERSION, server.BuildID,
		this.conf.IdgenWorkerId, time.Since(this.startedAt))

	profiler.do(IDENT, ctx, "{r^%s}", r)

	return
}
