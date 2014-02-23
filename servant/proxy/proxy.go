/*
Proxy of remote servant so that we can dispatch request
to cluster instead of having to serve all by ourselves.
*/
package proxy

import (
	"sync"
	"time"
)

type Proxy struct {
	mutex *sync.Mutex

	// pools for each remote peer(faed) instance
	capacity    int
	idleTimeout time.Duration
	pools       map[string]*funServantPeerPool
}

func New(capacity int, idleTimeout time.Duration) (this *Proxy) {
	this = &Proxy{capacity: capacity, idleTimeout: idleTimeout,
		mutex: new(sync.Mutex)}
	this.pools = make(map[string]*funServantPeerPool)
	return
}

func (this *Proxy) Servant(serverAddr string) (*FunServantPeer, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.getServant(serverAddr)
}

func (this *Proxy) getServant(serverAddr string) (*FunServantPeer, error) {
	if _, ok := this.pools[serverAddr]; !ok {
		this.pools[serverAddr] = newFunServantPeerPool(serverAddr,
			this.capacity, this.idleTimeout)
		this.pools[serverAddr].Open()
	}

	return this.pools[serverAddr].Get()
}
