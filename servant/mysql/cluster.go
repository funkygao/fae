package mysql

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
)

type MysqlCluster struct {
	conf     *config.ConfigMysql
	selector ServerSelector
	breakers map[string]*breaker.Consecutive
	clients  map[string]*mysql // key is pool name
}

func New(cf *config.ConfigMysql) *MysqlCluster {
	this := new(MysqlCluster)
	this.conf = cf
	this.breakers = make(map[string]*breaker.Consecutive)
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
