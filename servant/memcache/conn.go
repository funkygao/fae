package memcache

import (
	"bufio"
	"net"
	"time"
)

// conn is a connection to a server.
type conn struct {
	nc   net.Conn
	rw   *bufio.ReadWriter
	addr net.Addr

	client *Client
}

// release returns this connection back to the client's free pool
func (this *conn) release() {
	this.client.putFreeConn(this.addr, this)
}

// condRelease releases this connection if the error pointed to by err
// is is nil (not an error) or is only a protocol level error (e.g. a
// cache miss).  The purpose is to not recycle TCP connections that
// are bad.
func (this *conn) condRelease(err *error) {
	if *err == nil || resumableError(*err) {
		this.release()
	} else {
		this.nc.Close()
	}
}

func (this *conn) extendDeadline() {
	this.nc.SetDeadline(time.Now().Add(this.client.conf.Timeout))
}
