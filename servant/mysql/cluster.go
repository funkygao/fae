package mysql

import (
	"github.com/funkygao/fae/config"
)

type MysqlCluster struct {
	conf     *config.ConfigMysql
	selector ServerSelector
	clients  map[string]*mysql // key is pool name
}

func New(cf *config.ConfigMysql) *MysqlCluster {
	this := new(MysqlCluster)
	this.conf = cf
	switch cf.ShardStrategy {
	default:
		this.selector = newStandardServerSelector()
	}
	this.selector.SetServers(cf)
	this.clients = make(map[string]*mysql)
	for _, pool := range cf.Pools() {
		this.clients[pool] = newMysql(cf.Servers[pool].DSN())
	}
	return this
}

func (this *MysqlCluster) Query(pool string, table string, shardId int32,
	sql string, args []interface{}) (r [][]byte, err error) {
	return
}
