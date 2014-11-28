/*
Proxy of remote servant so that we can dispatch request
to cluster instead of having to serve all by ourselves.
*/
package proxy

import (
	"encoding/json"
	"sync"
	"time"
)

type Proxy struct {
	mutex *sync.Mutex

	capacity    int           // all fae peer share same capacity, weight TODO
	idleTimeout time.Duration // fae peer in pool idle timeout

	pools map[string]*funServantPeerPool // each fae peer has a pool, key is peerAddr
}

func New(capacity int, idleTimeout time.Duration) *Proxy {
	return &Proxy{
		capacity:    capacity,
		idleTimeout: idleTimeout,
		mutex:       new(sync.Mutex),
		pools:       make(map[string]*funServantPeerPool),
	}
}

func (this *Proxy) Servant(peerAddr string) (*FunServantPeer, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.pools[peerAddr]; !ok {
		this.pools[peerAddr] = newFunServantPeerPool(peerAddr,
			this.capacity, this.idleTimeout)
		this.pools[peerAddr].Open()
	}

	return this.pools[peerAddr].Get()
}

func (this *Proxy) StatsJSON() string {
	m := make(map[string]string)
	for addr, pool := range this.pools {
		m[addr] = pool.pool.StatsJSON()
	}

	pretty, _ := json.MarshalIndent(m, "", "    ")
	return string(pretty)
}

func (this *Proxy) StatsMap() map[string]string {
	m := make(map[string]string)
	for addr, pool := range this.pools {
		m[addr] = pool.pool.StatsJSON()
	}

	return m
}
