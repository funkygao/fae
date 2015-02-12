package mysql

import (
	sql_ "database/sql"
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
	"time"
)

type MysqlCluster struct {
	selector ServerSelector
}

func New(cf *config.ConfigMysql) *MysqlCluster {
	this := new(MysqlCluster)
	switch cf.ShardStrategy {
	case "standard":
		this.selector = newStandardServerSelector(cf)

	case "vbucket":
		panic("vbucket mysql sharding not implemented")

	default:
		panic("unknown mysql sharding type: " + cf.ShardStrategy)
	}

	return this
}

func (this *MysqlCluster) Conn(pool string, table string,
	hintId int) (*sql_.DB, error) {
	my, err := this.selector.PickServer(pool, table, hintId)
	if err != nil {
		return nil, err
	}

	return my.db, nil
}

func (this *MysqlCluster) QueryShards(pool string, table string, sql string,
	args []interface{}) (cols []string, rows [][]string, ex error) {
	rows = make([][]string, 0)
	var (
		rawRowValues []sql_.RawBytes
		scanArgs     []interface{}
		rowValues    []string
	)
	// TODO query in parallel
	for _, my := range this.selector.PoolServers(pool) {
		rs, err := my.Query(sql, args...)
		if err != nil {
			ex = err
			return
		}

		if len(cols) == 0 {
			// initialize the vars only once
			cols, ex = rs.Columns()
			if ex != nil {
				rs.Close()
				return
			}

			rawRowValues = make([]sql_.RawBytes, len(cols))
			scanArgs = make([]interface{}, len(cols))
			for i, _ := range cols {
				scanArgs[i] = &rawRowValues[i]
			}
		}

		for rs.Next() {
			if ex = rs.Scan(scanArgs...); ex != nil {
				rs.Close()
				return
			}

			rowValues = make([]string, len(cols))
			// TODO O(N), room for optimization, allow_nullable_columns
			for i, raw := range rawRowValues {
				if raw == nil {
					rowValues[i] = "NULL"
				} else {
					rowValues[i] = string(raw)
				}
			}

			rows = append(rows, rowValues)
		}

		rs.Close()
	}

	return
}

func (this *MysqlCluster) Query(pool string, table string, hintId int,
	sql string, args []interface{}) (*sql_.Rows, error) {
	my, err := this.selector.PickServer(pool, table, hintId)
	if err != nil {
		return nil, err
	}

	return my.Query(sql, args...)
}

func (this *MysqlCluster) Exec(pool string, table string, hintId int,
	sql string, args []interface{}) (afftectedRows int64,
	lastInsertId int64, err error) {
	this.selector.KickLookupCache(pool, hintId)

	my, err := this.selector.PickServer(pool, table, hintId)
	if err != nil {
		return 0, 0, err
	}

	return my.Exec(sql, args...)
}

func (this *MysqlCluster) Close() (err error) {
	for _, my := range this.selector.Servers() {
		if e := my.db.Close(); e != nil {
			err = e
		}
	}
	return
}

func (this *MysqlCluster) KickLookupCache(pool string, hintId int) {
	this.selector.KickLookupCache(pool, hintId)
}

func (this *MysqlCluster) Warmup() {
	var (
		err error
		t1  = time.Now()
	)

	for _, m := range this.selector.Servers() {
		err = m.Ping()
		if err != nil {
			log.Error("Warmup mysql: %s", err)
			break
		}

	}

	if err != nil {
		log.Error("Mysql failed to warmup within %s: %s",
			time.Since(t1), err)
	} else {
		log.Debug("Mysql warmup within %s: %+v",
			time.Since(t1), this.selector)
	}
}
