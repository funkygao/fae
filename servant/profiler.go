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
	t1 time.Time
}

func (this *profiler) do(name string, ctx *rpc.Context, format string,
	args ...interface{}) {
	elapsed := time.Since(this.t1)
	slow := elapsed.Seconds() > 3 // TODO config
	if !(slow || this.on) {
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

func (this *profiler) contextInfo(ctx *rpc.Context) (r contextInfo) {
	const (
		N         = 3
		SEPERATOR = "+"
	)
	p := strings.SplitN(ctx.Caller, SEPERATOR, N)
	if len(p) != N {
		return
	}

	r.ctx = ctx
	r.httpMethod, r.uri, r.seqId = p[0], p[1], p[2]

	return
}

func (this *profiler) truncatedStr(val string) string {
	if len(val) < config.Servants.ProfilerMaxBodySize {
		return val
	}

	return val[:config.Servants.ProfilerMaxBodySize] + "..."
}
