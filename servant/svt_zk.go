package servant

import (
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) ZkCreate(ctx *rpc.Context, path string,
	data string) (r bool, appErr error) {
	const IDENT = "zk.create"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	// TODO always persistent?
	if err = etclib.Create(path, data, 0); err == nil {
		r = true
	}

	profiler.do(IDENT, ctx, "{path^%s data^%s} {r^%v err^%v}",
		path, string(data), r, err)
	return
}

func (this *FunServantImpl) ZkChildren(ctx *rpc.Context,
	path string) (r []string, appErr error) {
	const IDENT = "zk.children"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)
	r, err = etclib.Children(path)

	profiler.do(IDENT, ctx, "{path^%s} {r^%+v err^%v}",
		path, r, err)
	return
}

func (this *FunServantImpl) ZkDel(ctx *rpc.Context,
	path string) (r bool, appErr error) {
	const IDENT = "zk.del"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)
	if err = etclib.Delete(path); err == nil {
		r = true
	}

	profiler.do(IDENT, ctx, "{path^%s} {r^%v err^%v}",
		path, r, err)
	return
}
