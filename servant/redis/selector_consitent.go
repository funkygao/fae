package redis

import (
	"github.com/funkygao/golib/hash"
)

type ConsistentServerSelector struct {
	nodes *hash.Map
}

func (this *ConsistentServerSelector) SetServers(servers ...string) error {
	if this.nodes == nil {
		this.nodes = hash.New(32, nil) // TODO config replica
	}

	this.nodes.Add(servers...)

	return nil
}

func (this *ConsistentServerSelector) PickServer(key string) (addr string) {
	return this.nodes.Get(key)
}
