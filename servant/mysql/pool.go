package mysql

import (
	"github.com/funkygao/fae/config"
)

type ClientPool struct {
	conf    *config.ConfigMysql
	clients map[string]*SqlDb
}

func New(cf *config.ConfigMysql) *ClientPool {
	this := new(ClientPool)
	this.conf = cf
	this.clients = make(map[string]*SqlDb)
	for _, pool := range cf.Pools() {
		this.clients[pool] = NewSqlDb("mysql", dsn, nil)
	}
	return this
}
