package servant

import (
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) ZkCreate(ctx *rpc.Context, path string,
	data string) (r bool, ex error) {
	const IDENT = "zk.create"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		this.stats.incErr()
		return
	}

	this.stats.inc(IDENT)

	// TODO always persistent?
	if ex = etclib.Create(path, data, 0); ex == nil {
		r = true
	} else {
		this.stats.incErr()
	}

	profiler.do(IDENT, ctx, "{path^%s data^%s} {r^%v err^%v}",
		path, string(data), r, ex)
	return
}

func (this *FunServantImpl) ZkChildren(ctx *rpc.Context,
	path string) (r []string, ex error) {
	const IDENT = "zk.children"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		this.stats.incErr()
		return
	}

	this.stats.inc(IDENT)
	r, ex = etclib.Children(path)
	if ex != nil {
		this.stats.incErr()
	}

	profiler.do(IDENT, ctx, "{path^%s} {r^%+v err^%v}",
		path, r, ex)
	return
}

func (this *FunServantImpl) ZkDel(ctx *rpc.Context,
	path string) (r bool, ex error) {
	const IDENT = "zk.del"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		this.stats.incErr()
		return
	}

	this.stats.inc(IDENT)
	if ex = etclib.Delete(path); ex == nil {
		r = true
	} else {
		this.stats.incErr()
	}

	profiler.do(IDENT, ctx, "{path^%s} {r^%v err^%v}",
		path, r, ex)
	return
}
