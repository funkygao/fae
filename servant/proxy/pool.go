package proxy

import (
	"github.com/funkygao/golib/pool"
	"time"
)

type funServantPeerPool struct {
	serverAddr  string
	capacity    int
	idleTimeout time.Duration
	pool        *pool.ResourcePool
}

func newFunServantPeerPool(serverAddr string, capacity int,
	idleTimeout time.Duration) (this *funServantPeerPool) {
	this = &funServantPeerPool{idleTimeout: idleTimeout, capacity: capacity,
		serverAddr: serverAddr}
	return
}

func (this *funServantPeerPool) Open() {
	factory := func() (pool.Resource, error) {
		client, err := connect(this.serverAddr)
		if err != nil {
			return nil, err
		}

		return &FunServantPeer{FunServantClient: client, pool: this}, nil
	}

	this.pool = pool.NewResourcePool(factory,
		this.capacity, this.capacity,
		this.idleTimeout)
}

func (this *funServantPeerPool) IsClosed() bool {
	return this.pool == nil
}

func (this *funServantPeerPool) Get() (*FunServantPeer, error) {
	fun, err := this.pool.Get()
	if err != nil {
		return nil, err
	}

	return fun.(*FunServantPeer), nil
}

func (this *funServantPeerPool) Put(conn *FunServantPeer) {
	this.pool.Put(conn)
}
