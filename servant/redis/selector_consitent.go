package redis

import (
	"github.com/funkygao/golib/hash"
	"net"
	"strings"
)

type ConsistentServerSelector struct {
	nodes *hash.Map
}

func (this *ConsistentServerSelector) SetServers(servers ...string) error {
	if this.nodes == nil {
		this.nodes = hash.New(2, nil)
	}

	this.nodes.Add(servers...)

	return nil
}

func (this *ConsistentServerSelector) PickServer(key string) (net.Addr, error) {
	server := this.nodes.Get(key)
	if strings.Contains(server, "/") {
		return net.ResolveUnixAddr("unix", server)
	}

	return net.ResolveTCPAddr("tcp", server)
}

func (this *ConsistentServerSelector) ServerList() (servers []net.Addr) {
	return
}
