// http://go-database-sql.org/surprises.html
// http://jmoiron.net/blog/built-in-interfaces/

package servant

import (
	"crypto/sha1"
	sql_ "database/sql"
	_json "encoding/json"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	json "github.com/funkygao/go-simplejson"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/mergemap"
	"strings"
)

func (this *FunServantImpl) MyQuery(ctx *rpc.Context, pool string, table string,
	hintId int64, sql string, args []string, cacheKey string) (r *rpc.MysqlResult,
	ex error) {
	const IDENT = "my.query"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	var (
		cacheKeyHash = cacheKey
		peer         string
		rows         int
	)

	if cacheKey != "" && this.conf.Mysql.CacheKeyHash {
		hashSum := sha1.Sum([]byte(cacheKey)) // sha1.Size
		cacheKeyHash = string(hashSum[:])
	}

	if cacheKeyHash == "" {
		r, ex = this.doMyQuery(IDENT, ctx, pool, table, hintId,
			sql, args, cacheKeyHash)
		rows = len(r.Rows)
		if r.RowsAffected > 0 {
			rows = int(r.RowsAffected)
		}
	} else {
		if ctx.IsSetSticky() && *ctx.Sticky {
			r, ex = this.doMyQuery(IDENT, ctx, pool, table, hintId,
				sql, args, cacheKeyHash)
			rows = len(r.Rows)
			if r.RowsAffected > 0 {
				rows = int(r.RowsAffected)
			}
		} else {
			svt, err := this.proxy.ServantByKey(cacheKey)
			if err != nil {
				ex = err
				return
			}

			if svt == nil {
				r, ex = this.doMyQuery(IDENT, ctx, pool, table, hintId,
					sql, args, cacheKeyHash)
				rows = len(r.Rows)
				if r.RowsAffected > 0 {
					rows = int(r.RowsAffected)
				}
			} else {
				// dispatch to peer
				peer = svt.Addr()
				svt.HijackContext(ctx)
				r, ex = svt.MyQuery(ctx, pool, table, hintId, sql, args, cacheKey)
				if ex != nil {
					svt.Close()
				} else {
					rows = len(r.Rows)
					if r.RowsAffected > 0 {
						rows = int(r.RowsAffected)
					}
				}

				svt.Recycle() // NEVER forget about this
			}
		}
	}

	if ex != nil {
		profiler.do(IDENT, ctx,
			"P=%s {cache^%s pool^%s table^%s id^%d sql^%s args^%+v} {err^%s}",
			peer, cacheKey, pool, table, hintId, sql, args, ex)
	} else {
		profiler.do(IDENT, ctx,
			"P=%s {cache^%s pool^%s table^%s id^%d sql^%s args^%+v} {rows^%d r^%+v}",
			peer, cacheKey, pool, table, hintId, sql, args, rows, *r)
	}

	return
}

func (this *FunServantImpl) MyEvict(ctx *rpc.Context,
	cacheKey string) (ex error) {
	const IDENT = "my.evict"

	this.stats.inc(IDENT)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	var peer string
	if ctx.IsSetSticky() && *ctx.Sticky {
		this.dbCacheStore.Del(cacheKey)
	} else {
		svt, err := this.proxy.ServantByKey(cacheKey)
		if err != nil {
			ex = err
			return
		}

		if svt == nil {
			this.dbCacheStore.Del(cacheKey)
		} else {
			peer = svt.Addr()
			svt.HijackContext(ctx)
			ex = svt.MyEvict(ctx, cacheKey)
			if ex != nil {
				svt.Close()
			}

			svt.Recycle()
		}
	}

	profiler.do(IDENT, ctx, "{key^%s} {p^%s}", cacheKey, peer)

	return
}

// If conflicts, jsonVal prevails
func (this *FunServantImpl) MyMerge(ctx *rpc.Context, pool string, table string,
	hintId int64, where string, key string, column string,
	jsonVal string) (r *rpc.MysqlMergeResult, ex error) {
	const IDENT = "my.merge"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	// find the column value from db
	// TODO keep in mem, needn't query db on each call
	querySql := "SELECT " + column + " FROM " + table + " WHERE " + where
	queryResult, err := this.doMyQuery(IDENT, ctx, pool, table, hintId,
		querySql, nil, "")
	if err != nil {
		ex = err
		log.Error("%s[%s]: %s", IDENT, querySql, err.Error())
		return
	}
	if len(queryResult.Rows) != 1 {
		ex = ErrMyMergeInvalidRow
		return
	}

	this.mysqlMergeMutexMap.Lock(key)
	defer this.mysqlMergeMutexMap.Unlock(key)

	// do the merge in mem
	j1, err := json.NewJson([]byte(queryResult.Rows[0][0]))
	if err != nil {
		ex = err
		return
	}
	j2, err := json.NewJson([]byte(jsonVal))
	if err != nil {
		ex = err
		return
	}

	var m1, m2 map[string]interface{}
	if m1, ex = j1.Map(); ex != nil {
		return
	}
	if m2, ex = j2.Map(); ex != nil {
		return
	}

	// TODO who wins if conflict on the same key
	merged := mergemap.Merge(m1, m2)

	// update db with merged value
	newVal, err := _json.Marshal(merged)
	if err != nil {
		ex = err
		return
	}

	updateSql := "UPDATE " + table + " SET " + column + "='" +
		string(newVal) + "' WHERE " + where
	_, err = this.doMyQuery(IDENT, ctx, pool, table, hintId, updateSql,
		nil, "")
	if err != nil {
		log.Error("%s[%s]: %s", IDENT, updateSql, err.Error())
		ex = err
		return
	}

	r = rpc.NewMysqlMergeResult()
	r.Ok = true
	r.NewVal = string(newVal)

	profiler.do(IDENT, ctx,
		"{key^%s pool^%s table^%s id^%d} {ok^%v val^%s}",
		key, pool, table, hintId, r.Ok, r.NewVal)
	return
}

// TODO ServantByKey(cacheKey)
func (this *FunServantImpl) doMyQuery(ident string, ctx *rpc.Context,
	pool string, table string, hintId int64, sql string,
	args []string, cacheKey string) (r *rpc.MysqlResult, ex error) {
	const (
		SQL_SELECT = "SELECT"
		SQL_UPDATE = "UPDATE"
	)

	// convert []string to []interface{}
	iargs := make([]interface{}, len(args), len(args))
	for i, arg := range args {
		iargs[i] = arg
	}

	r = rpc.NewMysqlResult()
	if strings.HasPrefix(sql, SQL_SELECT) { // SELECT MUST be in upper case
		ex = this.doMySelect(r, ident, ctx, pool, table, hintId,
			sql, args, iargs, cacheKey)
	} else {
		ex = this.doMyExec(r, ident, ctx, pool, table, hintId,
			sql, args, iargs, cacheKey)
	}

	return
}

func (this *FunServantImpl) doMySelect(r *rpc.MysqlResult,
	ident string, ctx *rpc.Context,
	pool string, table string, hintId int64, sql string,
	args []string, iargs []interface{}, cacheKey string) (ex error) {
	if cacheKey != "" {
		if cacheValue, present := this.dbCacheStore.Get(cacheKey); present {
			log.Debug("Q=%s cache[%s] hit", ident, cacheKey)
			this.dbCacheHits.Inc("hit", 1)
			*r = *(cacheValue.(*rpc.MysqlResult))
			return
		}
	}

	// cache miss, do real db query
	rows, err := this.my.Query(pool, table, int(hintId), sql, iargs)
	if err != nil {
		ex = err
		return
	}

	// recycle the underlying connection back to conn pool
	defer rows.Close()

	// pack the result
	cols, err := rows.Columns()
	if err != nil {
		ex = err
		log.Error("Q=%s %s[%s]: sql=%s args=(%v): %s",
			ident,
			pool, table,
			sql, args,
			ex)
		return
	} else {
		r.Cols = cols
		r.Rows = make([][]string, 0)
		for rows.Next() {
			rawRowValues := make([]sql_.RawBytes, len(cols))
			scanArgs := make([]interface{}, len(cols))
			for i, _ := range cols {
				scanArgs[i] = &rawRowValues[i]
			}
			if ex = rows.Scan(scanArgs...); ex != nil {
				log.Error("Q=%s %s[%s]: sql=%s args=(%v): %s",
					ident,
					pool, table,
					sql, args,
					ex)
				return
			}

			rowValues := make([]string, len(cols))
			for i, raw := range rawRowValues {
				if raw == nil {
					rowValues[i] = "NULL"
				} else {
					rowValues[i] = string(raw)
				}
			}

			r.Rows = append(r.Rows, rowValues)
		}

		// check for errors after we’re done iterating over the rows
		if ex = rows.Err(); ex != nil {
			log.Error("Q=%s %s[%s]: sql=%s args=(%v): %s",
				ident,
				pool, table,
				sql, args,
				ex)
			return
		}

		// query success, set cache: even when empty data returned
		if cacheKey != "" {
			this.dbCacheStore.Set(cacheKey, r)

			this.dbCacheHits.Inc("miss", 1)
			log.Debug("Q=%s cache[%s] miss", ident, cacheKey)
		}
	}

	return
}

func (this *FunServantImpl) doMyExec(r *rpc.MysqlResult,
	ident string, ctx *rpc.Context,
	pool string, table string, hintId int64, sql string,
	args []string, iargs []interface{}, cacheKey string) (err error) {
	if r.RowsAffected, r.LastInsertId, err = this.my.Exec(pool,
		table, int(hintId), sql, iargs); err != nil {
		log.Error("Q=%s %s[%s]: sql=%s args=(%v): %s",
			ident, pool, table, sql, args, err)
		return
	}

	// update success, del cache
	if cacheKey != "" {
		this.dbCacheStore.Del(cacheKey)

		this.dbCacheHits.Inc("kicked", 1)
		log.Debug("Q=%s cache[%s] kicked", ident, cacheKey)
	}

	return
}
