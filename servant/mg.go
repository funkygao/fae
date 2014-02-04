package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) MgInsert(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, doc []byte, options []byte) (r bool, intError error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		intError = err
		return
	}

	err = sess.DB().C(table).Insert(doc)
	if err == nil {
		r = true
	}
	sess.Recyle(&err)

	return
}

func (this *FunServantImpl) MgDelete(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte) (r bool, intError error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		intError = err
		return
	}

	err = sess.DB().C(table).Remove(query)
	if err == nil {
		r = true
	}
	sess.Recyle(&err) // reuse this session, we should never forget this
	return
}

func (this *FunServantImpl) MgFindOne(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte, fields []byte) (r []byte, intError error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		intError = err
		return
	}

	err = sess.DB().C(table).Find(query).One(&r)
	sess.Recyle(&err)

	return
}

func (this *FunServantImpl) MgFindAll(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte, fields []byte, limit []byte,
	orderBy []byte) (r []byte, intError error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		intError = err
		return
	}

	sess.Recyle(&err)

	return
}

func (this *FunServantImpl) MgUpdate(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte, change []byte) (r bool, intError error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		intError = err
		return
	}

	err = sess.DB().C(table).Update(query, change)
	if err == nil {
		r = true
	}
	sess.Recyle(&err)

	return
}

func (this *FunServantImpl) MgUpsert(ctx *rpc.ReqCtx, kind string, shardId int32,
	table string, query []byte, change []byte) (r bool, intError error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		intError = err
		return
	}

	_, err = sess.DB().C(table).Upsert(query, change)
	if err == nil {
		r = true
	}
	sess.Recyle(&err)

	return
}

func (this *FunServantImpl) MgFindAndModify(ctx *rpc.ReqCtx, kind string,
	shardId int32, table string, command []byte) (r []byte, intError error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		intError = err
		return
	}

	sess.Recyle(&err)

	return
}
