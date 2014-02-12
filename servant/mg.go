/*
mongodb doc:json bytes
*/
package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/mongo"
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo/bson"
)

func (this *FunServantImpl) MgInsert(ctx *rpc.Context,
	kind string, table string, shardId int32,
	doc []byte) (r bool, appErr error) {
	profiler := this.profiler()

	// get mongodb session
	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	// unmarsal inbound param
	// client json_encode, server json_decode into internal bson.M struct
	bsonDoc, err := this.unmarshalJson(doc)
	if err != nil {
		appErr = err
		return
	}

	// do insert and check error
	err = sess.DB().C(table).Insert(bsonDoc)
	if err != nil {
		// will not rais app error
		log.Error(err)
	} else {
		r = true
	}

	// recycle the mongodb session
	sess.Recyle(&err)

	profiler.do("mg.insert", ctx,
		"{kind^%s table^%s id^%d doc^%s} {%v}",
		kind, table, shardId,
		this.truncatedBytes(doc),
		r)

	return
}

func (this *FunServantImpl) MgInserts(ctx *rpc.Context,
	kind string, table string, shardId int32,
	docs [][]byte) (r bool, appErr error) {
	profiler := this.profiler()

	// get mongodb session
	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	// unmarsal inbound param
	// client json_encode, server json_decode into internal bson.M struct
	bsonDocs := make([]interface{}, len(docs))
	for i, doc := range docs {
		bsonDoc, err := this.unmarshalJson(doc)
		if err != nil {
			appErr = err
			return
		}

		bsonDocs[i] = bsonDoc
	}

	// do insert and check error
	err := sess.DB().C(table).Insert(bsonDocs...)
	if err != nil {
		// will not rais app error
		log.Error(err)
	} else {
		r = true
	}

	// recycle the mongodb session
	sess.Recyle(&err)

	profiler.do("mg.inserts", ctx,
		"{kind^%s table^%s id^%d docs^%d} {%v}",
		kind, table, shardId,
		len(docs),
		r)
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

	bsonQuery, err := this.unmarshalJson(query)
	if err != nil {
		appErr = err
		return
	}
	err = sess.DB().C(table).Remove(bsonQuery)
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

	bsonQuery, err := this.unmarshalJson(query)
	if err != nil {
		appErr = err
		return
	}
	var bsonFields bson.M
	if len(fields) > 5 {
		log.Info("ha %d", len(fields))
		bsonFields, err = this.unmarshalJson(fields)
		if err != nil {
			appErr = err
			return
		}
	}

	var result bson.M
	err = sess.DB().C(table).Find(bsonQuery).Select(bsonFields).One(&result)
	if err != nil {
		log.Error(err)
		appErr = err
		return
	}
	sess.Recyle(&err)

	r, _ = bson.Marshal(result)

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
	orderBy []string) (r [][]byte, appErr error) {
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	bsonQuery, err := this.unmarshalJson(query)
	if err != nil {
		appErr = err
		return
	}
	bsonFields, err := this.unmarshalJson(fields)
	if err != nil {
		appErr = err
		return
	}

	q := sess.DB().C(table).Find(bsonQuery).Select(bsonFields)
	if limit > 0 {
		q.Limit(int(limit))
	}
	if skip > 0 {
		q.Skip(int(skip))
	}
	q.Sort(orderBy...)
	var result []bson.M
	appErr = q.All(&result)
	if appErr == nil {
		r = make([][]byte, len(result))
		for i, v := range result {
			r[i], _ = bson.Marshal(v)
		}
	}

	sess.Recyle(&appErr)

	profiler.do("mg.findAll", ctx,
		"{kind^%s table^%s id^%d query%s fields^%s} {%d}",
		kind, table, shardId,
		this.truncatedBytes(query),
		this.truncatedBytes(fields),
		len(r))

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

	bsonQuery, err := this.unmarshalJson(query)
	if err != nil {
		appErr = err
		return
	}
	bsonChange, err := this.unmarshalJson(change)
	if err != nil {
		appErr = err
		return
	}

	appErr = sess.DB().C(table).Update(bsonQuery, bsonChange)
	if appErr == nil {
		r = true
	}
	sess.Recyle(&appErr)

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

	bsonQuery, err := this.unmarshalJson(query)
	if err != nil {
		appErr = err
		return
	}
	bsonChange, err := this.unmarshalJson(change)
	if err != nil {
		appErr = err
		return
	}

	_, appErr = sess.DB().C(table).Upsert(bsonQuery, bsonChange)
	if appErr == nil {
		r = true
	}
	sess.Recyle(&appErr)

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
	profiler := this.profiler()

	var sess *mongo.Session
	sess, appErr = this.mongoSession(kind, shardId)
	if appErr != nil {
		return
	}

	bsonQuery, err := this.unmarshalJson(query)
	if err != nil {
		appErr = err
		return
	}

	var r int
	r, appErr = sess.DB().C(table).Find(bsonQuery).Count()
	n = int32(r)

	sess.Recyle(&err)

	profiler.do("mg.count", ctx,
		"{kind^%s table^%s id^%d query^%s} {%d}",
		kind, table, shardId,
		this.truncatedBytes(query),
		n)

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
