package proxy

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/pool"
	"time"
)

// a conn pool to a fae endpoint
type funServantPeerPool struct {
	serverAddr  string
	capacity    int
	idleTimeout time.Duration
	pool        *pool.ResourcePool
}

// A kind of pool.Resource
type FunServantPeer struct {
	pool.Resource
	*rpc.FunServantClient

	pool *funServantPeerPool
}

func (this *FunServantPeer) Close() {
	this.Transport.Close()
}

func (this *FunServantPeer) Recycle() {
	if this.Transport.IsOpen() {
		this.pool.Put(this)
	} else {
		this.pool.Put(nil)
	}
}

func (this *FunServantPeer) NewContext() *rpc.Context {
	ctx := rpc.NewContext()
	ctx.Rid = "1"
	ctx.Reason = "proxy"
	return ctx
}

func newFunServantPeerPool(serverAddr string, capacity int,
	idleTimeout time.Duration) (this *funServantPeerPool) {
	this = &funServantPeerPool{idleTimeout: idleTimeout, capacity: capacity,
		serverAddr: serverAddr}
	return
}

func (this *funServantPeerPool) connect(serverAddr string) (*rpc.FunServantClient,
	error) {
	transportFactory := thrift.NewTBufferedTransportFactory(2 << 10)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocketTimeout(serverAddr, 0)
	if err != nil {
		return nil, err
	}

	useTransport := transportFactory.GetTransport(transport)
	client := rpc.NewFunServantClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		return nil, err
	}

	return client, nil
}

func (this *funServantPeerPool) Open() {
	factory := func() (pool.Resource, error) {
		client, err := this.connect(this.serverAddr)
		if err != nil {
			return nil, err
		}

		return &FunServantPeer{FunServantClient: client, pool: this}, nil
	}

	this.pool = pool.NewResourcePool(factory,
		this.capacity, this.capacity,
		this.idleTimeout)
}

func (this *funServantPeerPool) Get() (*FunServantPeer, error) {
	fun, err := this.pool.Get()
	if err != nil {
		return nil, err
	}

	return fun.(*FunServantPeer), nil
}

func (this *funServantPeerPool) Put(conn *FunServantPeer) {
	if !this.pool.IsClosed() {
		this.pool.Put(conn)
	}
}
