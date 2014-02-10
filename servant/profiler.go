package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/sampling"
	log "github.com/funkygao/log4go"
	"time"
)

type profiler struct {
	on bool
	t1 time.Time
}

func (this *FunServantImpl) profiler() profiler {
	p := profiler{on: false}
	p.on = sampling.SampleRateSatisfied(this.conf.ProfilerRate)
	p.t1 = time.Now()

	return p
}

func (this *profiler) do(name string, ctx *rpc.Context, format interface{}, args ...interface{}) {
	elapsed := time.Since(this.t1)
	if elapsed.Seconds() > 5.0 { // TODO config
		// slow response
		s := fmt.Sprintf("SLOW T=%s Q=%s X{%s} "+format, elapsed, name, this.contextInfo(ctx), args...)
		log.Warn(s)
	} else if this.on {
		s := fmt.Sprintf("T=%s Q=%s X{%s} "+format, elapsed, name, this.contextInfo(ctx), args...)
		log.Debug(s)
	}
}

func (this *FunServantImpl) truncatedBytes(val []byte) []byte {
	if len(val) < this.conf.ProfilerMaxAnswerSize {
		return val
	}

	return append(val[:this.conf.ProfilerMaxAnswerSize], []byte("..."))
}
