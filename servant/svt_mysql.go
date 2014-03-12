package servant

import (
	//sql_ "database/sql"
	"encoding/json"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
	"reflect"
)

func (this *FunServantImpl) MyQuery(ctx *rpc.Context, pool string, table string,
	shardId int32, sql string, args [][]byte) (r []byte, appErr error) {
	profiler := this.profiler()
	rows, err := this.my.Query(pool, table, int(shardId), sql, nil)
	if err != nil {
		appErr = err
		log.Error("my.query: %v", err)
	}
	defer rows.Close()

	// pack the result
	res := make(map[string]interface{})
	cols, err := rows.Columns()
	if err != nil {
		appErr = err
		log.Error("my.query: %v", err)
	} else {
		res["cols"] = cols
		vals := make([][]string, 0)
		for rows.Next() {
			rawRowValues := make([]interface{}, len(cols))
			scanArgs := make([]interface{}, len(cols))
			rowValues := make([]string, len(cols))
			for i, _ := range cols {
				scanArgs[i] = &rawRowValues[i]
			}
			err = rows.Scan(scanArgs...)
			if err != nil {
				appErr = err
				log.Error("my.query: %v", err)
			}
			for i, raw := range rawRowValues {
				log.Debug("%v", reflect.TypeOf(raw))
				if raw == nil {
					rowValues[i] = "NULL"
				} else {
					//rowValues[i] = string(raw)
				}
			}

			vals = append(vals, rowValues)
		}
		res["vals"] = vals
	}

	r, _ = json.Marshal(res)
	profiler.do("my.query", ctx,
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
