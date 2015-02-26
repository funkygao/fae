package proxy

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/pool"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/thrift/lib/go/thrift"
	"net"
	"sync/atomic"
)

// a conn pool to a single fae remote peer
type funServantPeerPool struct {
	peerAddr string
	myIp     string

	cf config.ConfigProxy

	pool *pool.ResourcePool

	nextServantId uint64 // each conn in this pool has an id
	txn           int64
}

func newFunServantPeerPool(myIp string, peerAddr string,
	cf config.ConfigProxy) (this *funServantPeerPool) {
	this = &funServantPeerPool{
		myIp:     myIp,
		peerAddr: peerAddr,
		cf:       cf,
	}
	return
}

func (this *funServantPeerPool) Open() {
	factory := func() (pool.Resource, error) {
		client, err := this.connect(this.peerAddr)
		if err != nil {
			return nil, err
		}

		id := atomic.AddUint64(&this.nextServantId, 1)
		log.Debug("peer[%s] connected txn:%d", this.peerAddr, id)

		return newFunServantPeer(id, this, client), nil
	}

	this.pool = pool.NewResourcePool("peer", factory,
		this.cf.PoolCapacity, this.cf.PoolCapacity, this.cf.IdleTimeout,
		this.cf.DiagnosticInterval, this.cf.BorrowMaxSeconds)
}

func (this *funServantPeerPool) Close() {
	this.pool.Close()
}

func (this *funServantPeerPool) Get() (*FunServantPeer, error) {
	fun, err := this.pool.Get()
	if err != nil {
		return nil, err
	}

	return fun.(*FunServantPeer), nil
}

// connect to remote servant peer
func (this *funServantPeerPool) connect(peerAddr string) (*rpc.FunServantClient,
	error) {
	transportFactory := thrift.NewTBufferedTransportFactory(this.cf.BufferSize)
	transport, err := thrift.NewTSocketTimeout(peerAddr, this.cf.IoTimeout)
	if err != nil {
		return nil, err
	}

	if err = transport.Open(); err != nil {
		log.Error("connect peer[%s]: %s", peerAddr, err)

		return nil, err
	}

	if tcpConn, ok := transport.Conn().(*net.TCPConn); ok {
		// nagle's only applies to client rather than server
		tcpConn.SetNoDelay(this.cf.TcpNoDelay)
	}

	client := rpc.NewFunServantClientFactory(
		transportFactory.GetTransport(transport),
		thrift.NewTBinaryProtocolFactoryDefault())

	return client, nil
}

func (this *funServantPeerPool) nextTxn() int64 {
	return atomic.AddInt64(&this.txn, 1)
}
