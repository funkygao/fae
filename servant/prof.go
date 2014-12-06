package servant

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"time"
)

// profiler and auditter
type profiler struct {
	on bool
	t0 time.Time // start of each session
	t1 time.Time // start of each call
}

func (this *profiler) do(callName string, ctx *rpc.Context, format string,
	args ...interface{}) {
	elapsed := time.Since(this.t1)
	slow := elapsed > config.Servants.CallSlowThreshold
	if !(slow || this.on) {
		return
	}

	body := fmt.Sprintf(format, args...)
	if slow {
		header := fmt.Sprintf("SLOW=%s/%s Q=%s %s ",
			elapsed, time.Since(this.t0), callName, ctx.String())
		log.Warn(header + this.truncatedStr(body))
	} else if this.on {
		header := fmt.Sprintf("T=%s/%s Q=%s %s ",
			elapsed, time.Since(this.t0), callName, ctx.String())
		log.Debug(header + this.truncatedStr(body))
	}

}

func (this *profiler) truncatedStr(val string) string {
	if len(val) < config.Servants.ProfilerMaxBodySize {
		return val
	}

	return val[:config.Servants.ProfilerMaxBodySize] + "..."
}
