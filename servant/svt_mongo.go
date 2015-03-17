package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/mongo"
	"github.com/funkygao/golib/trace"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/thrift/lib/go/thrift"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func (this *FunServantImpl) MgInsert(ctx *rpc.Context,
	pool string, table string, shardId int32,
	doc []byte) (r bool, ex error) {
	const IDENT = "mg.insert"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	// unmarshal inbound param
	// client json_encode, server json_decode into internal bson.M struct
	bsonDoc, err := mongo.UnmarshalIn(doc)
	if err != nil {
		ex = err
		profiler.do(IDENT, ctx,
			"{pool^%s table^%s doc^%v} {err^%v r^%v}",
			pool, table,
			bsonDoc,
			ex,
			r)

		return
	}

	// do insert and check error
	//raw := bson.Raw{}
	//raw.Kind = '\x03'
	//raw.Data = doc
	err = sess.DB().C(table).Insert(bsonDoc)
	if err != nil {
		// will not rais app error
		log.Error("mg.insert: %v", err)
	} else {
		r = true
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s doc^%v} {err^%v r^%v}",
		pool, table,
		bsonDoc,
		ex,
		r)

	return
}

func (this *FunServantImpl) MgInserts(ctx *rpc.Context,
	pool string, table string, shardId int32,
	docs [][]byte) (r bool, ex error) {
	const IDENT = "mg.inserts"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	// get mongodb session
	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	// unmarsal inbound param
	// client bson_encode, server bson_decode into internal bson.M struct
	bsonDocs := make([]interface{}, len(docs))
	for i, doc := range docs {
		bsonDoc, err := mongo.UnmarshalIn(doc)
		if err != nil {
			ex = err
			return
		}

		bsonDocs[i] = bsonDoc
	}

	// do insert and check error
	err = sess.DB().C(table).Insert(bsonDocs...)
	if err != nil {
		// will not rais app error
		log.Error("%s: %v", IDENT, err)
	} else {
		r = true
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s docN^%d} {err^%v r^%v}",
		pool, table,
		len(docs),
		ex,
		r)

	return
}

func (this *FunServantImpl) MgDelete(ctx *rpc.Context,
	pool string, table string, shardId int32,
	query []byte) (r bool, ex error) {
	const IDENT = "mg.del"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	// get mongodb session
	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	bsonQuery, err := mongo.UnmarshalIn(query)
	if err != nil {
		ex = err
		return
	}
	err = sess.DB().C(table).Remove(bsonQuery)
	if err == nil {
		r = true
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s query^%v} {err^%v r^%v}",
		pool, table,
		bsonQuery,
		ex,
		r)

	return
}

func (this *FunServantImpl) MgFindOne(ctx *rpc.Context,
	pool string, table string, shardId int32,
	query []byte, fields []byte) (r []byte,
	miss *rpc.TMongoNotFound, ex error) {
	const IDENT = "mg.findOne"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	// get mongodb session
	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	bsonQuery, err := mongo.UnmarshalIn(query)
	if err != nil {
		ex = err
		return
	}
	var bsonFields bson.M
	if !mongo.FieldsIsNil(fields) {
		bsonFields, err = mongo.UnmarshalIn(fields)
		if err != nil {
			ex = err
			return
		}
	}

	var result bson.M
	q := sess.DB().C(table).Find(bsonQuery)
	if !mongo.FieldsIsNil(fields) {
		q.Select(bsonFields)
	}
	err = q.One(&result)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Error("%s: %v", IDENT, err)
		} else {
			miss = rpc.NewTMongoNotFound()
			miss.Message = thrift.StringPtr(err.Error())
			profiler.do(IDENT, ctx,
				"{pool^%s table^%s query^%v fields^%v} {miss^%v err^%v val^%v}",
				pool, table,
				bsonQuery,
				bsonFields,
				miss,
				ex,
				result)
			return
		}

		ex = err
		return
	}

	r = mongo.MarshalOut(result)

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s query^%v fields^%v} {miss^%v err^%v val^%v}",
		pool, table,
		bsonQuery,
		bsonFields,
		miss,
		ex,
		result)

	return
}

func (this *FunServantImpl) MgFindAll(ctx *rpc.Context,
	pool string, table string, shardId int32,
	query []byte, fields []byte, limit int32, skip int32,
	orderBy []string) (r [][]byte, ex error) {
	const IDENT = "mg.findAll"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	bsonQuery, err := mongo.UnmarshalIn(query)
	if err != nil {
		ex = err
		return
	}
	var bsonFields bson.M
	if !mongo.FieldsIsNil(fields) {
		bsonFields, err = mongo.UnmarshalIn(fields)
		if err != nil {
			ex = err
			return
		}
	}

	q := sess.DB().C(table).Find(bsonQuery)
	if !mongo.FieldsIsNil(fields) {
		q.Select(bsonFields)
	}
	if limit > 0 {
		q.Limit(int(limit))
	}
	if skip > 0 {
		q.Skip(int(skip))
	}
	if len(orderBy) > 0 {
		q.Sort(orderBy...)
	}

	var result []bson.M
	err = q.All(&result)
	if err == nil {
		r = make([][]byte, len(result))
		for i, v := range result {
			r[i] = mongo.MarshalOut(v)
		}
	} else {
		ex = err
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s query^%v fields^%v} {err^%v rN^%d}",
		pool, table,
		bsonQuery,
		bsonFields,
		ex,
		len(r))

	return
}

func (this *FunServantImpl) MgUpdate(ctx *rpc.Context,
	pool string, table string, shardId int32,
	query []byte, change []byte) (r bool, ex error) {
	const IDENT = "mg.update"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	// get mongodb session
	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	bsonQuery, err := mongo.UnmarshalIn(query)
	if err != nil {
		ex = err
		return
	}
	bsonChange, err := mongo.UnmarshalIn(change)
	if err != nil {
		ex = err
		return
	}

	err = sess.DB().C(table).Update(bsonQuery, bsonChange)
	if err == nil {
		r = true
	} else {
		log.Error("%s: %v", IDENT, err)
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s query^%v chg^%v} {err^%v r^%v}",
		pool, table,
		bsonQuery,
		bsonChange,
		ex,
		r)

	return
}

func (this *FunServantImpl) MgUpdateId(ctx *rpc.Context,
	pool string, table string, shardId int32,
	id int32, change []byte) (r bool, ex error) {
	ex = ErrNotImplemented
	return
}

func (this *FunServantImpl) MgUpsert(ctx *rpc.Context,
	pool string, table string, shardId int32,
	query []byte, change []byte) (r bool, ex error) {
	const IDENT = "mg.upsert"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	bsonQuery, err := mongo.UnmarshalIn(query)
	if err != nil {
		ex = err
		return
	}
	bsonChange, err := mongo.UnmarshalIn(change)
	if err != nil {
		ex = err
		return
	}

	_, err = sess.DB().C(table).Upsert(bsonQuery, bsonChange)
	if err == nil {
		r = true
	} else {
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s query^%v chg^%v} {err^%v r^%v}",
		pool, table,
		bsonQuery,
		bsonChange,
		ex,
		r)

	return
}

func (this *FunServantImpl) MgUpsertId(ctx *rpc.Context,
	pool string, table string, shardId int32,
	id int32, change []byte) (r bool, ex error) {
	ex = ErrNotImplemented
	return
}

func (this *FunServantImpl) MgCount(ctx *rpc.Context,
	pool string, table string, shardId int32,
	query []byte) (n int32, ex error) {
	const IDENT = "mg.count"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	// get mongodb session
	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	bsonQuery, err := mongo.UnmarshalIn(query)
	if err != nil {
		ex = err
		return
	}

	var r int
	r, ex = sess.DB().C(table).Find(bsonQuery).Count()
	n = int32(r)
	if ex != nil {
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s query^%v} {err^%v r^%d}",
		pool, table,
		bsonQuery,
		ex,
		n)

	return
}

func (this *FunServantImpl) MgFindAndModify(ctx *rpc.Context,
	pool string, table string, shardId int32,
	query []byte, change []byte, upsert bool,
	remove bool, returnNew bool) (r []byte, ex error) {
	const IDENT = "mg.findAndModify"

	if this.mg == nil {
		ex = ErrServantNotStarted
		return
	}

	svtStats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	// get mongodb session
	sess, err := this.mongoSession(pool, shardId)
	if err != nil {
		ex = err
		return
	}
	defer sess.Recyle(&err)

	bsonQuery, err := mongo.UnmarshalIn(query)
	if err != nil {
		ex = err
		return
	}
	bsonChange, err := mongo.UnmarshalIn(change)
	if err != nil {
		ex = err
		return
	}

	doc := bson.M{}
	changeInfo, _ := sess.DB().C(table).Find(bsonQuery).
		Apply(mgo.Change{Update: bsonChange,
		Upsert: upsert, Remove: remove, ReturnNew: returnNew}, &doc)
	r = mongo.MarshalOut(doc)

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s query^%v chg^%v} {err^%v updated^%d removed^%d r^%v}",
		pool, table,
		bsonQuery,
		bsonChange,
		ex,
		changeInfo.Updated, changeInfo.Removed,
		doc)
	return
}

func (this *FunServantImpl) MgFindId(ctx *rpc.Context,
	pool string, table string, shardId int32,
	id []byte) (r []byte, ex error) {
	ex = ErrNotImplemented
	return
}

type mongoProtocolLogger struct {
}

func (this *mongoProtocolLogger) Output(calldepth int, s string) error {
	log.Debug("(%s) %s", trace.CallerFuncName(calldepth), s)
	return nil
}

func (this *FunServantImpl) mongoSession(pool string,
	shardId int32) (*mongo.Session, error) {
	sess, err := this.mg.Session(pool, shardId)
	if err != nil {
		log.Error("{pool^%s id^%d}: %s", pool, shardId, err)
		return nil, err
	}

	return sess, err
}
