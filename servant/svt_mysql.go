// http://go-database-sql.org/surprises.html
// http://jmoiron.net/blog/built-in-interfaces/

package servant

import (
	sql_ "database/sql"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"strings"
)

func (this *FunServantImpl) MyQuery(ctx *rpc.Context, pool string, table string,
	hintId int64, sql string, args []string) (r *rpc.MysqlResult, appErr error) {
	const (
		IDENT      = "my.query"
		SQL_SELECT = "SELECT"
	)

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	// convert []string to []interface{}
	margs := make([]interface{}, len(args), len(args))
	for i, arg := range args {
		margs[i] = arg
	}

	r = rpc.NewMysqlResult()
	if strings.HasPrefix(sql, SQL_SELECT) { // SELECT MUST be in upper case
		rows, err := this.my.Query(pool, table, int(hintId), sql, margs)
		if err != nil {
			appErr = err
			log.Error("Q=%s %s: %s (%v) %s", IDENT, ctx.String(), sql, args, appErr)
			return
		}

		// recycle the underlying connection back to conn pool
		defer rows.Close()

		// pack the result
		cols, err := rows.Columns()
		if err != nil {
			appErr = err
			log.Error("Q=%s %s: %s (%v) %s", IDENT, ctx.String(), sql, args, appErr)
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
					log.Error("Q=%s %s: %s (%v) %s", IDENT, ctx.String(), sql, args, appErr)
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
				log.Error("Q=%s %s: %s (%v) %s", IDENT, ctx.String(), sql, args, appErr)
				return
			}
		}
	} else {
		if r.RowsAffected, r.LastInsertId, appErr = this.my.Exec(pool,
			table, int(hintId), sql, margs); appErr != nil {
			log.Error("Q=%s %s: %s (%v) %s", IDENT, ctx.String(), sql, args, appErr)
			return
		}
	}

	profiler.do(IDENT, ctx,
		"{pool^%s table^%s id^%d sql^%s args^%+v} {r^%#v}",
		pool, table, hintId, sql, args, *r)
	return
}
