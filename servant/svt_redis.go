package servant

import (
	"encoding/json"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/redigo/redis"
	"strconv"
)

func (this *FunServantImpl) RdCall(ctx *rpc.Context, cmd string,
	pool string, keysAndArgs []string) (r string, ex error) {
	const IDENT = "rd.call"

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	var val interface{}
	// cannot use args (type []string) as type []interface {}
	iargs := make([]interface{}, len(keysAndArgs))
	for i, v := range keysAndArgs {
		iargs[i] = v
	}
	if val, ex = this.rd.Call(cmd, pool, iargs...); ex == nil && val != nil {
		switch val := val.(type) {
		case []byte:
			r = string(val)

		case string:
			r = val

		case int64:
			// e,g. hset
			r = strconv.FormatInt(val, 10)

		case []interface{}:
			// e,g. hgetall
			strs, err := redis.Strings(val, nil)
			if err == nil {
				bytes, err := json.Marshal(strs)
				if err == nil {
					r = string(bytes)
				}
			}

		default:
			log.Error("redis.%s unknown result type: %T", cmd, val)
		}
	}

	if ex != nil {
		svtStats.incErr()
	}

	profiler.do(IDENT, ctx,
		"{cmd^%s pool^%s key^%s args^%+v} {r^%s}",
		cmd, pool, keysAndArgs[0], keysAndArgs[1:], r)

	return
}
