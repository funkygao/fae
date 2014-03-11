package mysql

import (
	"github.com/funkygao/fae/config"
)

type ClientPool struct {
	selector ServerSelector
	conf     *config.ConfigMysql
	clients  map[string]*SqlDb
}

func New(cf *config.ConfigMysql) *ClientPool {
	this := new(ClientPool)
	this.conf = cf
	switch cf.ShardStrategy {
	default:
		this.selector = newStandardServerSelector()
	}
	this.selector.SetServers(cf)
	this.clients = make(map[string]*SqlDb)
	for _, pool := range cf.Pools() {
		this.clients[pool] = newSqlDb("mysql", cf.Servers[pool].DSN(), nil)
	}
	return this
}

func (this *ClientPool) Query(pool string, table string, shardId int32,
	sql string, args []interface{}) (r [][]byte, err error) {
	return
}
