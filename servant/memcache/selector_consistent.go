package memcache

import (
	"net"
	"strings"
	"sync"
)

// Not implemented yet TODO
type ConsistentServerSelector struct {
	lk    sync.RWMutex
	addrs []net.Addr
}

func (this *ConsistentServerSelector) SetServers(servers ...string) error {
	naddr := make([]net.Addr, len(servers))
	for i, server := range servers {
		if strings.Contains(server, "/") {
			addr, err := net.ResolveUnixAddr("unix", server)
			if err != nil {
				return err
			}
			naddr[i] = addr
		} else {
			tcpaddr, err := net.ResolveTCPAddr("tcp", server)
			if err != nil {
				return err
			}
			naddr[i] = tcpaddr
		}
	}

	this.lk.Lock()
	defer this.lk.Unlock()
	this.addrs = naddr
	return nil
}

func (this *ConsistentServerSelector) PickServer(key string) (net.Addr, error) {
	return nil, nil
}
