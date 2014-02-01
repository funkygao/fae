package memcache

import (
	"net"
)

type ServerSelector interface {
	SetServers(servers ...string) error
	PickServer(key string) (net.Addr, error)
}
