package mongo

import (
	"github.com/funkygao/fae/config"
	"strings"
)

// Sharding was done on client
type LegacyServerSelector struct {
	Servers map[string]*config.ConfigMongodbServer // key is pool
}

func NewLegacyServerSelector(baseNum int) *LegacyServerSelector {
	return &LegacyServerSelector{}
}

func (this *LegacyServerSelector) PickServer(pool string,
	shardId int) (server *config.ConfigMongodbServer, err error) {
	var present bool
	server, present = this.Servers[this.normalizedPool(pool)]
	if !present {
		err = ErrServerNotFound
	}

	return
}

func (this *LegacyServerSelector) SetServers(servers map[string]*config.ConfigMongodbServer) {
	this.Servers = servers
}

func (this *LegacyServerSelector) normalizedPool(pool string) string {
	const (
		N      = 2
		PREFIX = "database."
	)

	p := strings.SplitN(pool, PREFIX, N)
	if len(p) == 2 {
		return p[1]
	}

	return pool
}
