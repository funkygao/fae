package engine

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"time"
)

type TUnixSocket struct {
	listener      net.Listener
	addr          net.Addr
	clientTimeout time.Duration
	interrupted   bool
}

func NewTUnixSocket(listenAddr string) (*TUnixSocket, error) {
	return NewTUnixSocketTimeout(listenAddr, 0)
}

func NewTUnixSocketTimeout(listenAddr string,
	clientTimeout time.Duration) (*TUnixSocket, error) {
	addr, err := net.ResolveUnixAddr("unix", listenAddr)
	if err != nil {
		return nil, err
	}
	return &TUnixSocket{addr: addr, clientTimeout: clientTimeout}, nil
}

func (this *TUnixSocket) Listen() error {
	if this.IsListening() {
		return nil
	}
	l, err := net.Listen(this.addr.Network(), this.addr.String())
	if err != nil {
		return err
	}
	this.listener = l
	return nil
}

func (this *TUnixSocket) Accept() (thrift.TTransport, error) {
	if this.interrupted {
		return nil, errors.New("Transport Interrupted")
	}
	if this.listener == nil {
		return nil, thrift.NewTTransportException(thrift.NOT_OPEN,
			"No underlying server socket")
	}
	conn, err := this.listener.Accept()
	if err != nil {
		return nil, thrift.NewTTransportExceptionFromError(err)
	}
	return thrift.NewTSocketFromConnTimeout(conn, this.clientTimeout), nil
}

func (this *TUnixSocket) IsListening() bool {
	return this.listener != nil
}

func (this *TUnixSocket) Open() error {
	if this.IsListening() {
		return thrift.NewTTransportException(thrift.ALREADY_OPEN,
			"Server socket already open")
	}
	if l, err := net.Listen(this.addr.Network(), this.addr.String()); err != nil {
		return err
	} else {
		this.listener = l
	}
	return nil
}

func (this *TUnixSocket) Addr() net.Addr {
	return this.addr
}

func (this *TUnixSocket) Close() error {
	defer func() {
		this.listener = nil
	}()
	if this.IsListening() {
		return this.listener.Close()
	}
	return nil
}

func (this *TUnixSocket) Interrupt() error {
	this.interrupted = true
	return nil
}
