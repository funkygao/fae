package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/idgen"
	"math/rand"
	"time"
)

// Ticket service
func (this *FunServantImpl) IdNext(ctx *rpc.Context) (r int64, ex error) {
	const IDENT = "id.next"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	for i := 0; i < 3; i++ {
		r, ex = this.idgen.Next()
		if ex != nil {
			// encounter ntp clock backwards problem, just retry
			time.Sleep(time.Millisecond * time.Duration(1+rand.Int63n(50)))
		} else {
			// got it!
			break
		}

	}

	profiler.do(IDENT, ctx, "{r^%d}", r)

	return
}

func (this *FunServantImpl) IdNextWithTag(ctx *rpc.Context,
	tag int16) (r int64, ex error) {
	const IDENT = "id.nextag"

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	for i := 0; i < 3; i++ {
		r, ex = this.idgen.NextWithTag(tag)
		if ex != nil {
			// encounter ntp clock backwards problem, just retry
			time.Sleep(time.Millisecond * time.Duration(1+rand.Int63n(50)))
		} else {
			// got it!
			break
		}
	}

	profiler.do(IDENT, ctx, "{tag^%d} {r^%d}", tag, r)

	return
}

func (this *FunServantImpl) IdDecode(ctx *rpc.Context,
	id int64) (r []int64, ex error) {
	const IDENT = "id.decode"
	svtStats.inc(IDENT)
	ts, tag, wid, seq := idgen.DecodeId(id)
	r = []int64{ts, tag, wid, seq}
	return
}
