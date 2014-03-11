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

func (this *MysqlCluster) Query(pool string, shardId int,
	sql string, args []interface{}) (*sql.Rows, error) {
	my, err := this.selector.PickServer(pool, shardId)
	if err != nil {
		return nil, err
	}

	return my.Query(sql, args...)
}
