package mongo

import (
	"github.com/funkygao/fae/config"
)

type ServerSelector interface {
	SetServers(servers map[string]*config.ConfigMongodbServer)
	PickServer(kind string, shardId int) (string, error)
}
