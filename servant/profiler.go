package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/sampling"
	log "github.com/funkygao/log4go"
	"time"
)

// profiler and auditter
type profiler struct {
	*FunServantImpl
	on bool
	t1 time.Time
}

func (this *FunServantImpl) profiler() *profiler {
	p := &profiler{on: false}
	p.on = sampling.SampleRateSatisfied(this.conf.ProfilerRate) // rand(1000) <= ProfilerRate
	p.t1 = time.Now()
	p.FunServantImpl = this

	return p
}

func (this *profiler) do(name string, ctx *rpc.Context, format string,
	args ...interface{}) {
	elapsed := time.Since(this.t1)
	slow := elapsed.Seconds() > 3
	if !slow && !this.on {
		return
	}

	// format visible body
	for idx, arg := range args {
		if mcData, ok := arg.(*rpc.TMemcacheData); ok {
			if mcData != nil {
				args[idx] = fmt.Sprintf("{Data:%s Flags:%d}",
					mcData.Data, mcData.Flags)
			}
		}
	}

	body := fmt.Sprintf(format, args...)
	if slow { // TODO config
		// slow response
		header := fmt.Sprintf("SLOW=%-10s Q=%s X{%s} ",
			elapsed, name, this.contextInfo(ctx))
		log.Warn(header + this.truncatedStr(body))
	} else if this.on {
		header := fmt.Sprintf("T=%-10s Q=%s X{%s} ",
			elapsed, name, this.contextInfo(ctx))
		log.Debug(header + this.truncatedStr(body))
	}
}

func (this *profiler) truncatedStr(val string) string {
	if len(val) < this.conf.ProfilerMaxBodySize {
		return val
	}

	return val[:this.conf.ProfilerMaxBodySize] + "..."
}
