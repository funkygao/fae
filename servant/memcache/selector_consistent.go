package memcache

import (
	"github.com/funkygao/golib/hash"
	"net"
	"strings"
)

type ConsistentServerSelector struct {
	peers *hash.Map
}

func (this *ConsistentServerSelector) SetServers(servers ...string) error {
	if this.peers == nil {
		this.peers = hash.New(2, nil)
	}

	this.peers.Add(servers...)

	return nil
}

func (this *ConsistentServerSelector) PickServer(key string) (net.Addr, error) {
	server := this.peers.Get(key)
	if strings.Contains(server, "/") {
		return net.ResolveUnixAddr("unix", server)
	}

	return net.ResolveTCPAddr("tcp", server)
}
