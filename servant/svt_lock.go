package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) Lock(ctx *rpc.Context,
	reason string, key string) (r bool, ex error) {
	const IDENT = "lock"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	var peer string
	if ctx.IsSetSticky() && *ctx.Sticky {
		svtStats.incPeerCall()

		r = this.lk.Lock(key)
	} else {
		svt, err := this.proxy.ServantByKey(key) // FIXME add prefix?
		if err != nil {
			ex = err
			if svt != nil {
				if proxy.IsIoError(err) {
					svt.Close()
				}
				svt.Recycle()
			}
			return
		}

		if svt == proxy.Self {
			r = this.lk.Lock(key)
		} else {
			svtStats.incCallPeer()

			peer = svt.Addr()
			svt.HijackContext(ctx)
			r, ex = svt.Lock(ctx, reason, key)
			if ex != nil {
				if proxy.IsIoError(ex) {
					svt.Close()
				}
			}

			svt.Recycle()
		}
	}

	profiler.do(IDENT, ctx, "P=%s {reason^%s key^%s} {r^%v}",
		peer, reason, key, r)

	if !r {
		log.Warn("P=%s lock failed: {reason^%s key^%s}", peer, reason, key)
	}

	return
}

func (this *FunServantImpl) Unlock(ctx *rpc.Context,
	reason string, key string) (ex error) {
	const IDENT = "unlock"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	var peer string
	if ctx.IsSetSticky() && *ctx.Sticky {
		svtStats.incPeerCall()

		this.lk.Unlock(key)
	} else {
		svt, err := this.proxy.ServantByKey(key)
		if err != nil {
			ex = err
			if svt != nil {
				if proxy.IsIoError(err) {
					svt.Close()
				}
				svt.Recycle()
			}
			return
		}

		if svt == proxy.Self {
			this.lk.Unlock(key)
		} else {
			svtStats.incCallPeer()

			peer = svt.Addr()
			svt.HijackContext(ctx)
			ex = svt.Unlock(ctx, reason, key)
			if ex != nil {
				if proxy.IsIoError(ex) {
					svt.Close()
				}
			}

			svt.Recycle()
		}
	}

	profiler.do(IDENT, ctx, "P=%s {reason^%s key^%s}",
		peer, reason, key)
	return
}
