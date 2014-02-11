/*
mongodb doc:json bytes
*/
package servant

import (
	"encoding/json"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/mongo"
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo/bson"
)

func (this *FunServantImpl) MgInsert(ctx *rpc.Context,
	kind string, table string, shardId int32,
	doc []byte, options []byte) (r bool, intError error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, intError = this.mongoSession(kind, shardId)
	if intError != nil {
		return
	}

	var bdoc = bson.M{}
	json.Unmarshal(doc, &bdoc)

	err := sess.DB().C(table).Insert(bdoc)
	if err == nil {
		r = true
	} else {
		log.Error(err)
	}
	sess.Recyle(&err)

	profiler.do("mg.ins", ctx,
		"{kind^%s table^%s id^%d doc^%s option^%s} {%v}",
		kind, table, shardId,
		doc, options,
		r)

	return
}

func (this *FunServantImpl) MgDelete(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte) (r bool, intError error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, intError = this.mongoSession(kind, shardId)
	if intError != nil {
		return
	}

	err := sess.DB().C(table).Remove(query)
	if err == nil {
		r = true
	}
	sess.Recyle(&err) // reuse this session, we should never forget this

	profiler.do("mg.del", ctx,
		"{kind^%s table^%s id^%d} {%v}",
		kind, table, shardId, r)
	return
}

func (this *FunServantImpl) MgFindOne(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, fields []byte) (r []byte, intError error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, intError = this.mongoSession(kind, shardId)
	if intError != nil {
		return
	}

	err := sess.DB().C(table).Find(query).One(&r)
	sess.Recyle(&err)

	profiler.do("mg.findOne", ctx,
		"{kind^%s table^%s id^d query%s fields^%s} {%s}",
		kind, table, shardId,
		query, fields,
		this.truncatedBytes(r))

	return
}

func (this *FunServantImpl) MgFindAll(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, fields []byte, limit []byte,
	orderBy []byte) (r []byte, intError error) {

	return
}

func (this *FunServantImpl) MgUpdate(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, change []byte) (r bool, intError error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, intError = this.mongoSession(kind, shardId)
	if intError != nil {
		return
	}

	err := sess.DB().C(table).Update(query, change)
	if err == nil {
		r = true
	}
	sess.Recyle(&err)

	profiler.do("mg.update", ctx,
		"{kind^%s table^%s id^%d query^%s chg^%s} {%v}",
		kind, table, shardId,
		query, change,
		r)

	return
}

func (this *FunServantImpl) MgUpsert(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, change []byte) (r bool, intError error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, intError = this.mongoSession(kind, shardId)
	if intError != nil {
		return
	}

	_, err := sess.DB().C(table).Upsert(query, change)
	if err == nil {
		r = true
	}
	sess.Recyle(&err)

	profiler.do("mg.upsert", ctx,
		"{kind^%s table^%s id^%d query^%s chg^%s} {%v}",
		kind, table, shardId,
		query, change,
		r)

	return
}

func (this *FunServantImpl) MgFindAndModify(ctx *rpc.Context,
	kind string, table string, shardId int32,
	command []byte) (r []byte, intError error) {

	return
}

func (this *FunServantImpl) mongoSession(kind string,
	shardId int32) (*mongo.Session, error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return sess, err
}
