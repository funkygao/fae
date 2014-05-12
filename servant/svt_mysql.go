// http://go-database-sql.org/surprises.html
// http://jmoiron.net/blog/built-in-interfaces/

package servant

import (
	sql_ "database/sql"
	"encoding/json"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"strings"
)

func (this *FunServantImpl) isSelectQuery(sql string) bool {
	return strings.HasPrefix(strings.ToLower(sql), "select")
}

func (this *FunServantImpl) MyQuery(ctx *rpc.Context, pool string, table string,
	shardId int32, sql string, args [][]byte) (r *rpc.MysqlResult, appErr error) {
	const IDENT = "my.query"

	profiler := this.profiler()
	this.stats.inc(IDENT)
	this.stats.inBytes.Inc(int64(len(sql)))
	profiler.do(IDENT, ctx,
		"{pool^%s table^%s id^%d sql^%s}",
		pool, table, shardId, sql)

	r = rpc.NewMysqlResult()
	if this.isSelectQuery(sql) {
		rows, err := this.my.Query(pool, table, int(shardId), sql, args)
		if err != nil {
			appErr = err
			log.Error("%s: %s %v", IDENT, sql, err)
			return
		}
		// recycle the underlying connection back to conn pool
		defer rows.Close()

		// pack the result
		res := make(map[string]interface{})
		cols, err := rows.Columns()
		if err != nil {
			appErr = err
			log.Error("%s: %s %v", IDENT, sql, err)
			return
		} else {
			res["cols"] = cols
			vals := make([][]string, 0)
			for rows.Next() {
				rawRowValues := make([]sql_.RawBytes, len(cols))
				scanArgs := make([]interface{}, len(cols))
				for i, _ := range cols {
					scanArgs[i] = &rawRowValues[i]
				}
				err = rows.Scan(scanArgs...)
				if err != nil {
					appErr = err
					log.Error("%s: %v", IDENT, err)
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

				vals = append(vals, rowValues)
			}
			// check for errors after we’re done iterating over the rows
			err = rows.Err()
			if err != nil {
				appErr = err
				log.Error("%s: %v", IDENT, err)
			}
			res["vals"] = vals
		}

		r.Rows, _ = json.Marshal(res)
	} else {
		var err error
		r.RowsAffected, r.LastInsertId, err = this.my.Exec(pool, table, int(shardId), sql, nil)
		if err != nil {
			appErr = err
			log.Error("%s: %s %v", IDENT, sql, err)
			return
		}
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s sql^%s} {r^%s}",
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
