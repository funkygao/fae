package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// actor lock
// register an uniq player name
// place a new player into a random tile in kingdom map
// under maintenance

// get a uniq name with length 3
func (this *FunServantImpl) GmName3(ctx *rpc.Context) (r string, appErr error) {
	const IDENT = "gm.name3"

	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	r = this.namegen.Next()

	profiler.do(IDENT, ctx, "{r^%s}", r)

	return
}

// TODO add reason and histogram of bytes
func (this *FunServantImpl) GmLatency(ctx *rpc.Context, ms int32,
	bytes int32) (appErr error) {
	this.phpLatency.Update(int64(ms))
	log.Trace("php rid^%s: {%dms, %dB}", ctx.Rid, ms, bytes)
	return
}

func (this *FunServantImpl) gm_register(ctx *rpc.Context, udid string) (appErr error) {
	return nil
}

func (this *FunServantImpl) gm_actor_lockuser(ctx *rpc.Context, uid int64) (appErr error) {
	return nil
}

func (this *FunServantImpl) gm_actor_locktile(ctx *rpc.Context, geohash int64) (appErr error) {
	appErr = ErrNotImplemented
	return nil
}
