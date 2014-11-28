/*
Proxy of remote servant so that we can dispatch request
to cluster instead of having to serve all by ourselves.
*/
package proxy

import (
	"encoding/json"
	"github.com/funkygao/etclib"
	log "github.com/funkygao/log4go"
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
	this := &Proxy{
		capacity:    capacity,
		idleTimeout: idleTimeout,
		mutex:       new(sync.Mutex),
		pools:       make(map[string]*funServantPeerPool),
	}

	return this
}

func (this *Proxy) StartMonitorCluster() {
	this.loadClusterSnapshot()
	go this.watchClusterPeers()
}

func (this *Proxy) loadClusterSnapshot() {
	faeNodes, err := etclib.ClusterNodes(etclib.NODE_FAE)
	if err != nil {
		log.Error("loadSnapshot[%s]: %s", etclib.NODE_FAE, err)
		return
	}

	for _, peerAddr := range faeNodes {
		// TODO discard self fae node
		// peerAddr is like "12.3.11.2:9001"
		this.Servant(peerAddr)

		log.Info("Found fae peer: %s", peerAddr)
	}

	log.Debug("cluster snapshot: %+v", this.StatsMap())
}

func (this *Proxy) watchClusterPeers() {
	for evt := range etclib.WatchFaeNodes() {
		log.Trace("cluster evt: %+v", evt)

		// TODO if self evt, ignore

		this.mutex.Lock()
		switch evt.EventType {
		case etclib.NODE_EVT_BOOT:
			this.Servant(evt.Addr)

		case etclib.NODE_EVT_SHUTDOWN:
			delete(this.pools, evt.Addr)
		}

		this.mutex.Unlock()
	}

}

// Get or create a fae peer servant based on peer address
// NOT goroutine safe, must set lock
func (this *Proxy) Servant(peerAddr string) (*FunServantPeer, error) {
	if _, ok := this.pools[peerAddr]; !ok {
		this.pools[peerAddr] = newFunServantPeerPool(peerAddr,
			this.capacity, this.idleTimeout)
		this.pools[peerAddr].Open()
	}

	return this.pools[peerAddr].Get()
}

func (this *Proxy) ClusterServants() map[string]*FunServantPeer {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// TODO ignore self
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
