package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) MgDelete(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte) (r bool, intError error) {
	sess, _ := this.mg.Session(kind, shardId)
	err := sess.DB().C(table).Remove(query)
	sess.Recyle(&err) // reuse this session, we should never forget this
	return
}

func (this *FunServantImpl) MgFindOne(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte, fields []byte) (r []byte, intError error) {
	return
}

func (this *FunServantImpl) MgFindAll(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte, fields []byte, limit []byte,
	orderBy []byte) (r []byte, intError error) {
	return
}

func (this *FunServantImpl) MgUpdate(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte, data []byte, upsert bool) (r bool, intError error) {
	return
}

func (this *FunServantImpl) MgInsert(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, data []byte, options []byte) (r bool, intError error) {
	return
}

func (this *FunServantImpl) MgInserts(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, data []byte, options []byte) (r bool, intError error) {
	return
}

func (this *FunServantImpl) MgFindAndModify(ctx *rpc.ReqCtx, kind string,
	shardId int32, table string, command []byte) (r []byte, intError error) {
	return
}
