package servant

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"strings"
	"time"
)

// profiler and auditter
type profiler struct {
	on bool
	t0 time.Time // start of each session
	t1 time.Time // start of each call
}

func (this *profiler) do(name string, ctx *rpc.Context, format string,
	args ...interface{}) {
	elapsed := time.Since(this.t1)
	slow := elapsed > config.Servants.CallSlowThreshold
	if !(slow || this.on) {
		return
	}

	body := fmt.Sprintf(format, args...)
	if slow {
		header := fmt.Sprintf("SLOW=%s/%s Q=%s X{%s} ",
			elapsed, time.Since(this.t0), name, this.contextInfo(ctx))
		log.Warn(header + this.truncatedStr(body))
	} else if this.on {
		header := fmt.Sprintf("T=%s/%s Q=%s X{%s} ",
			elapsed, time.Since(this.t0), name, this.contextInfo(ctx))
		log.Debug(header + this.truncatedStr(body))
	}

	// reset t1
	this.t1 = time.Now()
}

func (this *profiler) contextInfo(ctx *rpc.Context) (r contextInfo) {
	const (
		N         = 3
		SEPERATOR = "+"
	)

	// TODO discard Caller
	p := strings.SplitN(ctx.Caller, SEPERATOR, N)
	if len(p) != N {
		return
	}

	r.ctx = ctx
	r.httpMethod, r.uri = p[0], p[1]

	return
}

func (this *profiler) truncatedStr(val string) string {
	if len(val) < config.Servants.ProfilerMaxBodySize {
		return val
	}

	return val[:config.Servants.ProfilerMaxBodySize] + "..."
}
