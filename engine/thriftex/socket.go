package thriftex

import (
	//log "code.google.com/p/log4go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"time"
)

// thrift.TTransport implementation
type TSocket struct {
	conn    net.Conn
	addr    net.Addr // remote address
	timeout time.Duration
}

func NewTSocketTimeout(hostPort string, timeout time.Duration) (*TSocket, error) {
	addr, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		return nil, err
	}

	return NewTSocketFromAddrTimeout(addr, timeout), nil
}

func NewTSocketFromAddrTimeout(addr net.Addr, timeout time.Duration) *TSocket {
	return &TSocket{addr: addr, timeout: timeout}
}

func NewTSocketFromConnTimeout(conn net.Conn, timeout time.Duration) *TSocket {
	return &TSocket{conn: conn, addr: conn.RemoteAddr(), timeout: timeout}
}

func (this *TSocket) SetTimeout(timeout time.Duration) error {
	this.timeout = timeout
	return nil
}

func (this *TSocket) pushDeadline(read, write bool) {
	var t time.Time
	if this.timeout > 0 {
		t = time.Now().Add(time.Duration(this.timeout))
	}
	if read && write {
		this.conn.SetDeadline(t)
	} else if read {
		this.conn.SetReadDeadline(t)
	} else if write {
		this.conn.SetWriteDeadline(t)
	}
}

// Connects the socket, creating a new socket object if necessary.
func (this *TSocket) Open() error {
	if this.IsOpen() {
		return thrift.NewTTransportException(thrift.ALREADY_OPEN,
			"Socket already connected.")
	}
	if this.addr == nil {
		return thrift.NewTTransportException(thrift.NOT_OPEN,
			"Cannot open nil address.")
	}
	if len(this.addr.Network()) == 0 {
		return thrift.NewTTransportException(thrift.NOT_OPEN,
			"Cannot open bad network name.")
	}
	if len(this.addr.String()) == 0 {
		return thrift.NewTTransportException(thrift.NOT_OPEN,
			"Cannot open bad address.")
	}

	var err error
	if this.conn, err = net.DialTimeout(this.addr.Network(),
		this.addr.String(), this.timeout); err != nil {
		return thrift.NewTTransportException(thrift.NOT_OPEN, err.Error())
	}
	return nil
}

func (this *TSocket) Conn() net.Conn {
	return this.conn
}

func (this *TSocket) IsOpen() bool {
	if this.conn == nil {
		return false
	}

	return true
}

func (this *TSocket) Close() error {
	if this.conn != nil {
		err := this.conn.Close()
		if err != nil {
			return err
		}

		this.conn = nil
	}

	return nil
}

func (this *TSocket) Read(buf []byte) (int, error) {
	if !this.IsOpen() {
		return 0, thrift.NewTTransportException(thrift.NOT_OPEN,
			"Connection not open")
	}

	this.pushDeadline(true, false)
	n, err := this.conn.Read(buf)
	return n, thrift.NewTTransportExceptionFromError(err)
}

func (this *TSocket) Write(buf []byte) (int, error) {
	if !this.IsOpen() {
		return 0, thrift.NewTTransportException(thrift.NOT_OPEN,
			"Connection not open")
	}

	this.pushDeadline(false, true)
	return this.conn.Write(buf)
}

func (this *TSocket) Peek() bool {
	return this.IsOpen()
}

func (this *TSocket) Flush() error {
	return nil
}

func (this *TSocket) Interrupt() error {
	if !this.IsOpen() {
		return nil
	}

	return this.conn.Close()
}
