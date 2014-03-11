package mysql

import (
	"database/sql"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
)

type MysqlCluster struct {
	conf     *config.ConfigMysql
	selector ServerSelector
	breakers map[string]*breaker.Consecutive // key is dsn
	clients  map[string]*mysql               // key is dsn
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
	for _, server := range cf.Servers {
		this.clients[server.DSN()] = newMysql(server.DSN())
	}

	return this
}

func (this *MysqlCluster) Query(pool string, shardId int32,
	sql string, args []interface{}) (r *sql.Rows, err error) {
	return
}
