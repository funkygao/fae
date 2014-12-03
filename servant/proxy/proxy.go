/*
Proxy of remote servant so that we can dispatch request
to cluster instead of having to serve all by ourselves.
*/
package proxy

import (
	"encoding/json"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
	"sync"
)

type Proxy struct {
	mutex sync.Mutex
	cf    config.ConfigProxy
	pools map[string]*funServantPeerPool // each fae peer has a pool, key is peerAddr
}

func New(cf config.ConfigProxy) *Proxy {
	this := &Proxy{
		cf:    cf,
		pools: make(map[string]*funServantPeerPool),
	}

	return this
}

func (this *Proxy) StartMonitorCluster() {
	peersChan := make(chan []string, 10)
	go etclib.WatchService(etclib.SERVICE_FAE, peersChan)

	for {
		select {
		case peers := <-peersChan:
			log.Trace("Cluster peers event: %+v", peers)

			newpeers := this.recreatePeers(peers)
			log.Debug("Cluster lastest peers: %+v", newpeers)
		}
	}

}

func (this *Proxy) recreatePeers(peers []string) []string {
	this.mutex.Lock()

	for addr, _ := range this.pools {
		if addr == this.cf.SelfAddr {
			continue
		}

		delete(this.pools, addr)
	}

	newpeers := make([]string, 0)
	for _, addr := range peers {
		if addr == this.cf.SelfAddr {
			continue
		}

		this.Servant(addr)
		newpeers = append(newpeers, addr)
	}

	this.mutex.Unlock()

	return newpeers
}

// Get or create a fae peer servant based on peer address
func (this *Proxy) Servant(peerAddr string) (*FunServantPeer, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.pools[peerAddr]; !ok {
		this.pools[peerAddr] = newFunServantPeerPool(peerAddr,
			this.cf.PoolCapacity, this.cf.IdleTimeout)
		this.pools[peerAddr].Open()
	}

	return this.pools[peerAddr].Get()
}

// get all other servants in the cluster
// FIXME lock, but can't dead lock with this.Servant()
func (this *Proxy) ClusterServants() map[string]*FunServantPeer {
	rv := make(map[string]*FunServantPeer)
	for peerAddr, _ := range this.pools {
		svt, err := this.Servant(peerAddr)
		if err != nil {
			log.Error("peer servant[%s]: %s", peerAddr, err)
			continue
		}

		rv[peerAddr] = svt
	}

	return rv
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
