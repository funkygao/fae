package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/golib/trie"
	log "github.com/funkygao/log4go"
)

// get a uniq name with length 3
// TODO dump to redis periodically
func (this *FunServantImpl) GmName3(ctx *rpc.Context) (r string, appErr error) {
	const IDENT = "gm.name3"

	this.stats.inc(IDENT)
	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	if ctx.IsSetSticky() && *ctx.Sticky {
		// I' the final servant, got call from remote peers
		r = this.namegen.Next()
	} else {
		svt, err := this.proxy.ServantByKey(IDENT)
		if err != nil {
			appErr = err
			log.Error("%s: %s", IDENT, err)
			return
		}

		if svt == nil {
			// handle it by myself, got call locally
			r = this.namegen.Next()
		} else {
			// remote peer servant
			svt.HijackContext(ctx)
			r, appErr = svt.GmName3(ctx)
			if appErr != nil {
				log.Error("%s: %s", IDENT, appErr)

				svt.Close()
			} else {
				svt.Recycle()
			}
		}
	}

	profiler.do(IDENT, ctx, "{r^%s}", r)

	return
}

// record php request time and payload size in bytes
func (this *FunServantImpl) GmLatency(ctx *rpc.Context, ms int32,
	bytes int32) (appErr error) {
	this.phpLatency.Update(int64(ms))
	this.phpPayloadSize.Update(int64(bytes))

	log.Trace("{%dms %s}: {uid^%d rid^%s reason^%s}",
		ms, gofmt.ByteSize(bytes),
		this.extractUid(ctx), ctx.Rid, ctx.Reason)

	return
}

func (this *FunServantImpl) GmLike(ctx *rpc.Context,
	name string, mode int8) (r []string, appErr error) {
	t := trie.NewTrie() // TODO
	switch mode {
	case 1:
		r = t.PrefixSearch(name)

	case 2:
		r = t.FuzzySearch(name)
	}

	return
}
