package mongo

import (
	"github.com/funkygao/fae/config"
	"strings"
)

// Sharding was done on client
type LegacyServerSelector struct {
	servers map[string]*config.ConfigMongodbServer // key is pool
}

func NewLegacyServerSelector(baseNum int) *LegacyServerSelector {
	return &LegacyServerSelector{}
}

func (this *LegacyServerSelector) PickServer(pool string,
	shardId int) (server *config.ConfigMongodbServer, err error) {
	var present bool
	server, present = this.servers[this.normalizedPool(pool)]
	if !present {
		err = ErrServerNotFound
	}

	return
}

func (this *LegacyServerSelector) SetServers(servers map[string]*config.ConfigMongodbServer) {
	this.servers = servers
}

func (this *LegacyServerSelector) ServerList() (servers []*config.ConfigMongodbServer) {
	for _, s := range this.servers {
		servers = append(servers, s)
	}
	return
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
