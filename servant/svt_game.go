package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/golib/trie"
	log "github.com/funkygao/log4go"
)

// TODO use hset to reduce mem usage
// http://instagram-engineering.tumblr.com/post/12202313862/storing-hundreds-of-millions-of-simple-key-value
func (this *FunServantImpl) GmReserve(ctx *rpc.Context,
	tag string, oldName, newName string) (r bool, ex error) {
	const IDENT = "gm.reserve"
	const REDIS_POOL = "naming"
	var prefix = "acc:" + tag + ":"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	var counter string
	counter, ex = this.callRedis("INCR", REDIS_POOL, []string{prefix + newName})
	if counter == "1" {
		r = true

		if oldName != "" {
			_, ex = this.callRedis("DEL", REDIS_POOL, []string{prefix + oldName})
		}
	}

	profiler.do(IDENT, ctx, "{tag^%s old^%s new^%s} {r^%+v}",
		tag, oldName, newName, r)

	return
}

func (this *FunServantImpl) GmRegister(ctx *rpc.Context, typ string) (r int64,
	ex error) {
	const IDENT = "gm.reg"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	r, ex = this.game.Register(typ)

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
		return
	}

	var peer string
	if ctx.IsSetSticky() && *ctx.Sticky {
		// I' the final servant, got call from remote peers
		svtStats.incPeerCall()

		if !this.game.NameDbLoaded {
			this.game.NameDbLoaded = true
			this.loadName3Bitmap(ctx)
		}

		r = this.game.NextName()
	} else {
		svt, err := this.proxy.ServantByKey(IDENT)
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

		if svt == nil {
			// handle it by myself, got call locally
			if !this.game.NameDbLoaded {
				this.game.NameDbLoaded = true
				this.loadName3Bitmap(ctx)
			}

			r = this.game.NextName()
		} else {
			// remote peer servant
			svtStats.incCallPeer()

			peer = svt.Addr()
			svt.HijackContext(ctx)
			r, ex = svt.GmName3(ctx)
			if ex != nil {
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

	this.game.CheckIn(ctx.Uid)

	log.Trace("{%dms %s}: {uid^%d rid^%d reason^%s}",
		ms, gofmt.ByteSize(bytes), ctx.Uid, ctx.Rid, ctx.Reason)

	return
}

func (this FunServantImpl) GmPresence(ctx *rpc.Context,
	uids []int64) (r []bool, ex error) {
	const IDENT = "gm.presence"

	svtStats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	// first, get my instance online status
	r = this.game.OnlineStatus(uids)

	// then, get online status in the whole cluster(self excluded)
	if !ctx.IsSetSticky() {
		remoteSvts, err := this.proxy.RemoteServants(true)
		if err != nil {
			ex = err
			for _, svt := range remoteSvts {
				svt.Recycle()
			}
			return
		}

		for _, svt := range remoteSvts {
			svt.HijackContext(ctx) // don't bounce the call back again

			onlines, err := svt.GmPresence(ctx, uids)
			if err != nil {
				log.Error("%s: %s", IDENT, err)
				if proxy.IsIoError(err) {
					svt.Close()
				}
				svt.Recycle()
				continue // skip the remote err
			}

			for i, online := range onlines {
				if online {
					// in the cluster, if any peer vote for a user online, he's online
					r[i] = true
				}
			}

			svt.Recycle()
		}
	}

	profiler.do(IDENT, ctx, "{uids^%v} {r^%v}", uids, r)
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
