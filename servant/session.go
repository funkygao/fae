package servant

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/sampling"
	log "github.com/funkygao/log4go"
	"time"
)

type session struct {
	profiler *profiler
}

func (this *session) startProfiler() *profiler {
	if this.profiler == nil {
		this.profiler = &profiler{}
		// TODO 某些web server需要100%采样
		this.profiler.on = sampling.SampleRateSatisfied(config.Servants.ProfilerRate) // rand(1000) <= ProfilerRate
		this.profiler.t0 = time.Now()
		this.profiler.t1 = this.profiler.t0
	}

	this.profiler.t1 = time.Now()
	return this.profiler
}

func (this *FunServantImpl) getSession(ctx *rpc.Context) *session {
	s, present := this.sessions.Get(ctx.Rid)
	if !present {
		s = &session{}
		this.sessions.Set(ctx.Rid, s)

		log.Trace("new session {reason:%s rid:%d}", ctx.Reason, ctx.Rid)
	}

	return s.(*session)
}
