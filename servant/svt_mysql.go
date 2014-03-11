package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) MyQuery(ctx *rpc.Context, pool string, table string,
	shardId int32, sql string, args [][]byte) (r [][]byte, appErr error) {
	profiler := this.profiler()
	log.Debug("ha")
	profiler.do("my.query", ctx,
		"{pool^%s table^%s sql^%s} {r^%v}",
		pool, table, sql, r)
	return
}

func (this *FunServantImpl) MyQueryOne(ctx *rpc.Context, pool string, table string,
	shardId int32, sql string, args [][]byte) (r []byte, appErr error) {
	profiler := this.profiler()
	profiler.do("my.queryOne", ctx,
		"{pool^%s table^%s sql^%s} {r^%v}",
		pool, table, sql, r)
	return
}
