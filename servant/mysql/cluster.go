package mysql

import (
	"database/sql"
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
	default:
		this.selector = newStandardServerSelector(cf)
	}

	return this
}

func (this *MysqlCluster) Query(pool string, table string, hintId int,
	sql string, args []interface{}) (*sql.Rows, error) {
	my, err := this.selector.PickServer(pool, table, hintId)
	if err != nil {
		return nil, err
	}

	return my.Query(sql, args...)
}

func (this *MysqlCluster) Exec(pool string, table string, hintId int,
	sql string, args []interface{}) (afftectedRows int64,
	lastInsertId int64, err error) {
	my, err := this.selector.PickServer(pool, table, hintId)
	if err != nil {
		return 0, 0, err
	}

	return my.ExecSql(sql, args...)
}

func (this *MysqlCluster) Warmup() {
	var (
		err error
		t1  = time.Now()
	)

	for _, m := range this.selector.Servers() {
		err = m.Ping()
		if err != nil {
			log.Error(err)
			continue
		}

	}

	log.Trace("Mysql pool warmup finished within %s: %+v",
		time.Since(t1), this.selector)

}
