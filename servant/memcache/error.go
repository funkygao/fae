package memcache

import (
	"errors"
	"net"
)

var (
	ErrCacheMiss    = errors.New("memcache: cache miss")
	ErrCASConflict  = errors.New("memcache: compare-and-swap conflict")
	ErrNotStored    = errors.New("memcache: item not stored")
	ErrServerError  = errors.New("memcache: server error")
	ErrNoStats      = errors.New("memcache: no statistics available")
	ErrMalformedKey = errors.New("malformed: key is too long or contains invalid characters")
	ErrNoServers    = errors.New("memcache: no servers configured or available")
	ErrCircuitOpen  = errors.New("memcache: circuit open")
	ErrInvalidPool  = errors.New("memcache: invalid pool name")
)

// ConnectTimeoutError is the error type used when it takes
// too long to connect to the desired host. This level of
// detail can generally be ignored.
type ConnectTimeoutError struct {
	Addr net.Addr
}

func (this *ConnectTimeoutError) Error() string {
	return "memcache: connect timeout to " + this.Addr.String()
}
