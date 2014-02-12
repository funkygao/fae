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
	doc []byte, options []byte) (r bool, appErr error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
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

	profiler.do("mg.insert", ctx,
		"{kind^%s table^%s id^%d doc^%s option^%s} {%v}",
		kind, table, shardId,
		this.truncatedBytes(doc),
		this.truncatedBytes(options),
		r)

	return
}

func (this *FunServantImpl) MgInserts(ctx *rpc.Context,
	kind string, table string, shardId int32,
	doc [][]byte, options []byte) (r bool, appErr error) {
	return
}

func (this *FunServantImpl) MgDelete(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte) (r bool, appErr error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	var bquery = bson.M{}
	json.Unmarshal(query, &bquery)
	err := sess.DB().C(table).Remove(bquery)
	if err == nil {
		r = true
	}
	sess.Recyle(&err) // reuse this session, we should never forget this

	profiler.do("mg.del", ctx,
		"{kind^%s table^%s id^%d query^%s} {%v}",
		kind, table, shardId,
		this.truncatedBytes(query),
		r)
	return
}

func (this *FunServantImpl) MgFindOne(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, fields []byte) (r []byte, appErr error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	// TODO fields
	var bquery = bson.M{}
	json.Unmarshal(query, &bquery)
	err := sess.DB().C(table).Find(query).One(&r)
	sess.Recyle(&err)

	profiler.do("mg.findOne", ctx,
		"{kind^%s table^%s id^%d query^%s fields^%s} {%s}",
		kind, table, shardId,
		this.truncatedBytes(query), this.truncatedBytes(fields),
		this.truncatedBytes(r))

	return
}

func (this *FunServantImpl) MgFindAll(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, fields []byte, limit int32, skip int32,
	orderBy []string) (r []byte, appErr error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	var bquery = bson.M{}
	json.Unmarshal(query, &bquery)
	err := sess.DB().C(table).Find(query).All(&r)
	sess.Recyle(&err)

	profiler.do("mg.findAll", ctx,
		"{kind^%s table^%s id^%d query%s fields^%s} {%s}",
		kind, table, shardId,
		this.truncatedBytes(query), this.truncatedBytes(fields),
		this.truncatedBytes(r))

	return
}

func (this *FunServantImpl) MgUpdate(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, change []byte) (r bool, appErr error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	var bquery = bson.M{}
	json.Unmarshal(query, &bquery)
	var bchange = bson.M{}
	json.Unmarshal(change, &bchange)
	err := sess.DB().C(table).Update(bquery, bchange)
	if err == nil {
		r = true
	}
	sess.Recyle(&err)

	profiler.do("mg.update", ctx,
		"{kind^%s table^%s id^%d query^%s chg^%s} {%v}",
		kind, table, shardId,
		this.truncatedBytes(query),
		this.truncatedBytes(change),
		r)

	return
}

func (this *FunServantImpl) MgUpdateId(ctx *rpc.Context,
	kind string, table string, shardId int32,
	id int32, change []byte) (r bool, appErr error) {
	appErr = ErrNotImplemented
	return
}

func (this *FunServantImpl) MgUpsert(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte, change []byte) (r bool, appErr error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	var bquery = bson.M{}
	json.Unmarshal(query, &bquery)
	var bchange = bson.M{}
	json.Unmarshal(change, &bchange)
	_, err := sess.DB().C(table).Upsert(bquery, bchange)
	if err == nil {
		r = true
	}
	sess.Recyle(&err)

	profiler.do("mg.upsert", ctx,
		"{kind^%s table^%s id^%d query^%s chg^%s} {%v}",
		kind, table, shardId,
		this.truncatedBytes(query),
		this.truncatedBytes(change),
		r)

	return
}

func (this *FunServantImpl) MgUpsertId(ctx *rpc.Context,
	kind string, table string, shardId int32,
	id int32, change []byte) (r bool, appErr error) {
	appErr = ErrNotImplemented
	return
}

func (this *FunServantImpl) MgCount(ctx *rpc.Context,
	kind string, table string, shardId int32,
	query []byte) (n int32, appErr error) {
	return
}

func (this *FunServantImpl) MgFindAndModify(ctx *rpc.Context,
	kind string, table string, shardId int32,
	command []byte) (r []byte, appErr error) {

	return
}

func (this *FunServantImpl) MgFindId(ctx *rpc.Context,
	kind string, table string, shardId int32,
	id []byte) (r []byte, appErr error) {
	appErr = ErrNotImplemented
	return
}
