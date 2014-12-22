package proxy

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/pool"
	log "github.com/funkygao/log4go"
	"sync/atomic"
	"time"
)

// a conn pool to a single fae remote peer
type funServantPeerPool struct {
	peerAddr string

	capacity    int
	idleTimeout time.Duration
	pool        *pool.ResourcePool

	// ctx related
	txn  int64
	myIp string
}

func newFunServantPeerPool(myIp string, peerAddr string, capacity int,
	idleTimeout time.Duration) (this *funServantPeerPool) {
	this = &funServantPeerPool{idleTimeout: idleTimeout, capacity: capacity,
		peerAddr: peerAddr, myIp: myIp}
	return
}

func (this *funServantPeerPool) Open() {
	factory := func() (pool.Resource, error) {
		client, err := this.connect(this.peerAddr)
		if err != nil {
			return nil, err
		}

		log.Trace("peer[%s] connected", this.peerAddr)

		return newFunServantPeer(this, client), nil
	}

	this.pool = pool.NewResourcePool("FaePeer",
		factory,
		this.capacity, this.capacity,
		this.idleTimeout)
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

func (this *funServantPeerPool) connect(peerAddr string) (*rpc.FunServantClient,
	error) {
	transportFactory := thrift.NewTBufferedTransportFactory(4 << 10) // TODO
	transport, err := thrift.NewTSocket(peerAddr)                    // TODO, no timeout
	if err != nil {
		return nil, err
	}

	if err = transport.Open(); err != nil {
		log.Error("conn peer[%s]: %s", peerAddr, err)

		return nil, err
	}

	client := rpc.NewFunServantClientFactory(
		transportFactory.GetTransport(transport),
		thrift.NewTBinaryProtocolFactoryDefault())

	return client, nil
}

func (this *funServantPeerPool) nextTxn() int64 {
	return atomic.AddInt64(&this.txn, 1)
}
