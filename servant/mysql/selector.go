package mysql

import (
	"github.com/funkygao/fae/config"
)

type ServerSelector interface {
	SetServers(*config.ConfigMysql)
	PickServer(pool string, shardId int) (addr string, err error)
}
