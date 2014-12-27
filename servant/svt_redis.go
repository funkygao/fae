package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) RdCall(ctx *rpc.Context, cmd string,
	pool string, key string, args []string) (r string, appErr error) {
	const IDENT = "rd.call"

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	var val interface{}
	// cannot use args (type []string) as type []interface {}
	iargs := make([]interface{}, len(args))
	for i, v := range args {
		iargs[i] = v
	}
	if val, appErr = this.rd.Call(cmd, pool, key, iargs...); appErr == nil && val != nil {
		switch val := val.(type) {
		case []byte:
			r = string(val)

		case string:
			r = val
		}
	}

	profiler.do(IDENT, ctx,
		"{cmd^%s pool^%s key^%s arg^%+v} {r^%s}",
		cmd, pool, key, args, r)

	return
}
