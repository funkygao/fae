package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/sampling"
	log "github.com/funkygao/log4go"
	"time"
)

type profiler struct {
	*FunServantImpl
	on bool
	t1 time.Time
}

func (this *FunServantImpl) profiler() profiler {
	p := profiler{on: false}
	p.on = sampling.SampleRateSatisfied(this.conf.ProfilerRate)
	p.t1 = time.Now()
	p.FunServantImpl = this

	return p
}

func (this *profiler) do(name string, ctx *rpc.Context, format string,
	args ...interface{}) {
	elapsed := time.Since(this.t1)
	if elapsed.Seconds() > 5.0 { // TODO config
		// slow response
		body := fmt.Sprintf(format, args...)
		header := fmt.Sprintf("SLOW=%s Q=%s X{%s} ",
			elapsed, name, this.contextInfo(ctx))
		log.Warn(header + this.truncatedStr(body))
	} else if this.on {
		body := fmt.Sprintf(format, args...)
		header := fmt.Sprintf("T=%s Q=%s X{%s} ",
			elapsed, name, this.contextInfo(ctx))
		log.Debug(header + this.truncatedStr(body))
	}
}

func (this *profiler) truncatedStr(val string) string {
	if len(val) < this.conf.ProfilerMaxAnswerSize {
		return val
	}

	return val[:this.conf.ProfilerMaxAnswerSize] + "..."
}
