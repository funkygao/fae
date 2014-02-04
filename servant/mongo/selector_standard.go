package mongo

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"strings"
)

type StandardServerSelector struct {
	ShardBaseNum int
	Servers      map[string]*config.ConfigMongodbServer // key is kind
}

func NewStandardServerSelector(baseNum int) *StandardServerSelector {
	return &StandardServerSelector{ShardBaseNum: baseNum}
}

func (this *StandardServerSelector) PickServer(kind string,
	shardId int) (addr string, err error) {
	const SHARD_KIND_PREFIX = "db"

	var bucket string
	if !strings.HasPrefix(kind, SHARD_KIND_PREFIX) {
		bucket = kind
	} else {
		bucket = fmt.Sprintf("db%d", (shardId/this.ShardBaseNum)+1)
	}

	if server, present := this.Servers[bucket]; present {
		addr = server.Address()
	} else {
		err = ErrServerNotFound
	}

	return
}

func (this *StandardServerSelector) SetServers(servers map[string]*config.ConfigMongodbServer) {
	this.Servers = servers
}
