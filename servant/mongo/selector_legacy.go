package mongo

import (
	"github.com/funkygao/fae/config"
)

// Sharding was done on client
type LegacyServerSelector struct {
	Servers map[string]*config.ConfigMongodbServer // key is kind
}

func NewLegacyServerSelector(baseNum int) *LegacyServerSelector {
	return &LegacyServerSelector{}
}

func (this *LegacyServerSelector) PickServer(kind string,
	shardId int) (server *config.ConfigMongodbServer, err error) {
	var present bool
	server, present = this.Servers[kind]
	if !present {
		err = ErrServerNotFound
	}

	return
}

func (this *LegacyServerSelector) SetServers(servers map[string]*config.ConfigMongodbServer) {
	this.Servers = servers
}
