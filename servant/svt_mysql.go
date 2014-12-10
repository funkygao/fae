// http://go-database-sql.org/surprises.html
// http://jmoiron.net/blog/built-in-interfaces/

package servant

import (
	sql_ "database/sql"
	_json "encoding/json"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	json "github.com/funkygao/go-simplejson"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/mergemap"
	"strings"
)

func (this *FunServantImpl) MyQuery(ctx *rpc.Context, pool string, table string,
	hintId int64, sql string, args []string) (r *rpc.MysqlResult, appErr error) {
	const IDENT = "my.query"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	_, r, appErr = this.doMyQuery(IDENT, ctx, pool, table, hintId, sql, args)

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s id^%d sql^%s args^%+v} {r^%#v}",
		pool, table, hintId, sql, args, *r)
	return
}

// If conflicts, jsonVal prevails
func (this *FunServantImpl) MyMerge(ctx *rpc.Context, pool string, table string,
	hintId int64, where string, key string, column string,
	jsonVal string) (r *rpc.MysqlMergeResult, appErr error) {
	const IDENT = "my.merge"

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	// find the column value from db
	querySql := "SELECT " + column + " FROM " + table + " WHERE " + where
	_, queryResult, err := this.doMyQuery(IDENT, ctx, pool, table, hintId,
		querySql, nil)
	if err != nil {
		appErr = err
		log.Error("%s[%s]: %s", IDENT, querySql, err.Error())
		return
	}
	if len(queryResult.Rows) != 1 {
		appErr = ErrMyMergeInvalidRow
		return
	}

	this.lockmap.Lock(key)
	defer this.lockmap.Unlock(key)

	// do the merge in mem
	j1, err := json.NewJson([]byte(queryResult.Rows[0][0]))
	if err != nil {
		appErr = err
		return
	}
	j2, err := json.NewJson([]byte(jsonVal))
	if err != nil {
		appErr = err
		return
	}

	var m1, m2 map[string]interface{}
	if m1, appErr = j1.Map(); appErr != nil {
		return
	}
	if m2, appErr = j2.Map(); appErr != nil {
		return
	}

	// TODO who wins if conflict on the same key
	merged := mergemap.Merge(m1, m2)

	// update db with merged value
	newVal, err := _json.Marshal(merged)
	if err != nil {
		appErr = err
		return
	}

	updateSql := "UPDATE " + table + " SET " + column + "='" +
		string(newVal) + "' WHERE " + where
	_, _, err = this.doMyQuery(IDENT, ctx, pool, table, hintId, updateSql, nil)
	if err != nil {
		log.Error("%s[%s]: %s", IDENT, updateSql, err.Error())
		appErr = err
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

func (this *FunServantImpl) doMyQuery(ident string, ctx *rpc.Context,
	pool string, table string, hintId int64, sql string,
	args []string) (operation string, r *rpc.MysqlResult, appErr error) {
	const (
		SQL_SELECT = "SELECT"
		OP_QUERY   = "qry"
		OP_EXEC    = "exc"
	)

	// convert []string to []interface{}
	margs := make([]interface{}, len(args), len(args))
	for i, arg := range args {
		margs[i] = arg
	}

	r = rpc.NewMysqlResult()
	if strings.HasPrefix(sql, SQL_SELECT) { // SELECT MUST be in upper case
		operation = OP_QUERY

		rows, err := this.my.Query(pool, table, int(hintId), sql, margs)
		if err != nil {
			appErr = err
			log.Error("Q=%s %s %s[%s]: sql=%s args=(%v) %s", ident,
				ctx.String(),
				pool, table,
				sql, args,
				appErr)
			return
		}

		// recycle the underlying connection back to conn pool
		defer rows.Close()

		// pack the result
		cols, err := rows.Columns()
		if err != nil {
			appErr = err
			log.Error("Q=%s %s %s[%s]: sql=%s args=(%v) %s", ident,
				ctx.String(),
				pool, table,
				sql, args,
				appErr)
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
				if appErr = rows.Scan(scanArgs...); appErr != nil {
					log.Error("Q=%s %s %s[%s]: sql=%s args=(%v) %s", ident,
						ctx.String(),
						pool, table,
						sql, args,
						appErr)
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
			if appErr = rows.Err(); appErr != nil {
				log.Error("Q=%s %s %s[%s]: sql=%s args=(%v) %s", ident,
					ctx.String(),
					pool, table,
					sql, args,
					appErr)
				return
			}
		}
	} else {
		operation = OP_EXEC

		// FIXME if sql is 'select * from UesrInfo', runtime will get here
		if r.RowsAffected, r.LastInsertId, appErr = this.my.Exec(pool,
			table, int(hintId), sql, margs); appErr != nil {
			log.Error("Q=%s %s %s[%s]: sql=%s args=(%v) %s", ident,
				ctx.String(),
				pool, table,
				sql, args,
				appErr)
			return
		}
	}

	return
}
