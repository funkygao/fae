package memcache

import (
	"hash/crc32"
	"net"
	"strings"
	"sync"
)

type StandardServerSelector struct {
	lk    sync.RWMutex
	addrs []net.Addr
}

// Can have dup server address for higher weight
func (this *StandardServerSelector) SetServers(servers ...string) error {
	if len(servers) == 0 {
		return ErrNoServers
	}
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

func (this *StandardServerSelector) PickServer(key string) (net.Addr, error) {
	this.lk.RLock()
	defer this.lk.RUnlock()
	if len(this.addrs) == 0 {
		return nil, ErrNoServers
	}

	bucket := ((crc32.ChecksumIEEE([]byte(key)) >> 16) & 0x7fff) % uint32(len(this.addrs))
	return this.addrs[bucket], nil
}
