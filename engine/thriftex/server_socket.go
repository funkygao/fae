package thriftex

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"time"
)

// thrift.TTransport implementation
type TServerSocket struct {
	listener      net.Listener
	addr          net.Addr
	clientTimeout time.Duration
	interrupted   bool
}

func NewTServerSocket(listenAddr string) (*TServerSocket, error) {
	return NewTServerSocketTimeout(listenAddr, 0)
}

func NewTServerSocketTimeout(listenAddr string, clientTimeout time.Duration) (*TServerSocket, error) {
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	return &TServerSocket{addr: addr, clientTimeout: clientTimeout}, nil
}

func (p *TServerSocket) Listen() error {
	if p.IsListening() {
		return nil
	}
	l, err := net.Listen(p.addr.Network(), p.addr.String())
	if err != nil {
		return err
	}
	p.listener = l
	return nil
}

func (p *TServerSocket) Accept() (thrift.TTransport, error) {
	if p.interrupted {
		return nil, errTransportInterrupted
	}
	if p.listener == nil {
		return nil, thrift.NewTTransportException(thrift.NOT_OPEN,
			"No underlying server socket")
	}

	conn, err := p.listener.Accept()
	if err != nil {
		return nil, thrift.NewTTransportExceptionFromError(err)
	}

	return NewTSocketFromConnTimeout(conn, p.clientTimeout), nil
}

// Checks whether the socket is listening.
func (p *TServerSocket) IsListening() bool {
	return p.listener != nil
}

// Connects the socket, creating a new socket object if necessary.
func (p *TServerSocket) Open() error {
	if p.IsListening() {
		return thrift.NewTTransportException(thrift.ALREADY_OPEN,
			"Server socket already open")
	}

	if l, err := net.Listen(p.addr.Network(), p.addr.String()); err != nil {
		return err
	} else {
		p.listener = l
	}

	return nil
}

func (p *TServerSocket) Addr() net.Addr {
	return p.addr
}

func (p *TServerSocket) Close() error {
	defer func() {
		p.listener = nil
	}()
	if p.IsListening() {
		return p.listener.Close()
	}
	return nil
}

func (p *TServerSocket) Interrupt() error {
	p.interrupted = true
	return nil
}
