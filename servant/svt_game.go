package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/golib/trie"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) GmRegister(ctx *rpc.Context, typ string) (r int64,
	ex error) {
	const IDENT = "gm.reg"
	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	r, ex = this.game.Register(typ)
	if ex != nil {
		svtStats.incErr()
	}

	profiler.do(IDENT, ctx, "{type^%s} {r^%+v}", typ, r)

	return
}

// get a uniq name with length 3
func (this *FunServantImpl) GmName3(ctx *rpc.Context) (r string, ex error) {
	const IDENT = "gm.name3"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	var peer string
	if ctx.IsSetSticky() && *ctx.Sticky {
		// I' the final servant, got call from remote peers
		svtStats.incPeerCall()

		if !this.game.NameDbLoaded {
			this.game.NameDbLoaded = true
			go this.loadName3Bitmap(ctx)
		}

		r = this.game.NextName()
	} else {
		svt, err := this.proxy.ServantByKey(IDENT)
		if err != nil {
			ex = err
			svtStats.incErr()
			return
		}

		if svt == nil {
			// handle it by myself, got call locally
			if !this.game.NameDbLoaded {
				this.game.NameDbLoaded = true
				go this.loadName3Bitmap(ctx)
			}

			r = this.game.NextName()
		} else {
			// remote peer servant
			svtStats.incCallPeer()

			peer = svt.Addr()
			svt.HijackContext(ctx)
			r, ex = svt.GmName3(ctx)
			if ex != nil {
				svtStats.incErr()

				if proxy.IsIoError(ex) {
					svt.Close()
				}
			}

			svt.Recycle() // NEVER forget about this
		}
	}

	profiler.do(IDENT, ctx, "P=%s {r^%s}", peer, r)

	return
}

func (this *FunServantImpl) loadName3Bitmap(ctx *rpc.Context) {
	log.Trace("namegen snapshot loading...")

	result, err := this.doMyQuery("loadName3Bitmap", ctx,
		"ShardLookup", "AllianceLookup", 0,
		"SELECT acronym FROM AllianceLookup", nil, "")
	if err != nil {
		log.Error("namegen load snapshot: %s", err)
	} else {
		for _, row := range result.Rows {
			this.game.SetNameBusy(row[0])
		}
	}

	log.Trace("namegen snapshot loaded: %d rows", len(result.Rows))
}

// record php request time and payload size in bytes
func (this *FunServantImpl) GmLatency(ctx *rpc.Context, ms int32,
	bytes int32) (ex error) {
	const IDENT = "gm.latency"
	svtStats.inc(IDENT)

	this.game.UpdatePhpLatency(int64(ms))
	this.game.UpdatePhpPayloadSize(int64(bytes))

	log.Trace("{%dms %s}: {uid^%d rid^%s reason^%s}",
		ms, gofmt.ByteSize(bytes),
		this.extractUid(ctx), ctx.Rid, ctx.Reason)

	return
}

func (this *FunServantImpl) GmLock(ctx *rpc.Context,
	reason string, key string) (r bool, ex error) {
	const IDENT = "gm.lock"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	var peer string
	if ctx.IsSetSticky() && *ctx.Sticky {
		svtStats.incPeerCall()

		r = this.game.Lock(key)
	} else {
		svt, err := this.proxy.ServantByKey(key) // FIXME add prefix?
		if err != nil {
			ex = err
			svtStats.incErr()
			return
		}

		if svt == nil {
			r = this.game.Lock(key)
		} else {
			svtStats.incCallPeer()

			peer = svt.Addr()
			svt.HijackContext(ctx)
			r, ex = svt.GmLock(ctx, reason, key)
			if ex != nil {
				svtStats.incErr()

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

func (this *FunServantImpl) GmUnlock(ctx *rpc.Context,
	reason string, key string) (ex error) {
	const IDENT = "gm.unlock"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		svtStats.incErr()
		return
	}

	var peer string
	if ctx.IsSetSticky() && *ctx.Sticky {
		svtStats.incPeerCall()

		this.game.Unlock(key)
	} else {
		svt, err := this.proxy.ServantByKey(key)
		if err != nil {
			ex = err
			svtStats.incErr()
			return
		}

		if svt == nil {
			this.game.Unlock(key)
		} else {
			svtStats.incCallPeer()

			peer = svt.Addr()
			svt.HijackContext(ctx)
			ex = svt.GmUnlock(ctx, reason, key)
			if ex != nil {
				svtStats.incErr()

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

func (this *FunServantImpl) GmLike(ctx *rpc.Context,
	name string, mode int8) (r []string, ex error) {
	t := trie.NewTrie() // TODO
	switch mode {
	case 1:
		r = t.PrefixSearch(name)

	case 2:
		r = t.FuzzySearch(name)
	}

	return
}
