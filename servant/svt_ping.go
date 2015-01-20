package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/server"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, ex error) {
	const IDENT = "ping"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		this.stats.incErr()
		return
	}

	this.stats.inc(IDENT)

	r = fmt.Sprintf("pong, %s, myid:%d", server.BuildID, this.conf.IdgenWorkerId)

	profiler.do(IDENT, ctx, "{pong, %s, myid:%d}",
		server.BuildID, this.conf.IdgenWorkerId)

	return
}
