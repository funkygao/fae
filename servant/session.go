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
	ctx      *rpc.Context
}

func (this *session) startProfiler() (*profiler, error) {
	if this.profiler == nil {
		if this.ctx.Rid == "" || this.ctx.Reason == "" {
			log.Error("Invalid context: %s", this.ctx.String())
			return nil, ErrInvalidContext
		}

		this.profiler = &profiler{}
		// TODO 某些web server需要100%采样
		this.profiler.on = sampling.SampleRateSatisfied(config.Servants.ProfilerRate) // rand(1000) <= ProfilerRate
		this.profiler.t0 = time.Now()
		this.profiler.t1 = this.profiler.t0
	}

	this.profiler.t1 = time.Now()
	return this.profiler, nil
}

func (this *FunServantImpl) getSession(ctx *rpc.Context) *session {
	s, present := this.sessions.Get(ctx.Rid)
	if !present {
		s = &session{ctx: ctx}
		this.sessions.Set(ctx.Rid, s)

		log.Trace("new session {rid^%s reason^%s}", ctx.Rid, ctx.Reason)
	}

	return s.(*session)
}
