package mysql

import (
	"database/sql"
	"github.com/funkygao/fae/config"
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
