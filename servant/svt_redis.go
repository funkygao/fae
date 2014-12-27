package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) RdCall(ctx *rpc.Context, cmd string,
	pool string, keysAndArgs []string) (r string, appErr error) {
	const IDENT = "rd.call"

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	var val interface{}
	// cannot use args (type []string) as type []interface {}
	iargs := make([]interface{}, len(keysAndArgs))
	for i, v := range keysAndArgs {
		iargs[i] = v
	}
	if val, appErr = this.rd.Call(cmd, pool, iargs...); appErr == nil && val != nil {
		switch val := val.(type) {
		case []byte:
			r = string(val)

		case string:
			r = val
		}
	}

	profiler.do(IDENT, ctx,
		"{cmd^%s pool^%s key^%s args^%+v} {r^%s}",
		cmd, pool, keysAndArgs[0], keysAndArgs[1:], r)

	return
}
