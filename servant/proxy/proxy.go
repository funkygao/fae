package proxy

import (
	"encoding/json"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/ip"
	log "github.com/funkygao/log4go"
	"sync"
)

// Proxy of remote servant so that we can dispatch request
// to cluster instead of having to serve all by ourselves.
type Proxy struct {
	mutex sync.Mutex
	cf    config.ConfigProxy
	myIp  string

	remotePeerPools map[string]*funServantPeerPool // key is peerAddr
	selector  PeerSelector
}

func New(cf config.ConfigProxy) *Proxy {
	this := &Proxy{
		cf:        cf,
		remotePeerPools: make(map[string]*funServantPeerPool),
		selector:  newStandardPeerSelector(),
		myIp:      ip.LocalIpv4Addrs()[0],
	}

	return this
}

func (this *Proxy) Enabled() bool {
	return this.cf.Enabled()
}

func (this *Proxy) StartMonitorCluster() {
	if !this.Enabled() {
		log.Warn("servant proxy disabled")
		return
	}

	peersChan := make(chan []string, 10)
	go etclib.WatchService(etclib.SERVICE_FAE, peersChan)

	for {
		select {
		case <-peersChan:
			peers, err := etclib.ServiceEndpoints(etclib.SERVICE_FAE)
			if err == nil {
				// no lock, because running within 1 goroutine
				log.Trace("Cluster latest fae nodes: %+v", peers)

				this.selector.SetPeersAddr(this.recreatePeers(peers))
			} else {
				log.Error("Cluster peers: %s", err)
			}
		}
	}

	log.Warn("Cluster monitor died")
}

func (this *Proxy) recreatePeers(peers []string) []string {
	for peerAddr, _ := range this.remotePeerPools {
		if peerAddr == this.cf.SelfAddr {
			continue
		}

		this.remotePeerPools[peerAddr].Close()

		delete(this.remotePeerPools, peerAddr)
	}

	newpeers := make([]string, 0)
	for _, addr := range peers {
		if addr == this.cf.SelfAddr {
			continue
		}

		this.Servant(addr)
		newpeers = append(newpeers, addr)
	}

	return newpeers
}

// Get or create a fae peer servant based on peer address
func (this *Proxy) Servant(peerAddr string) (*FunServantPeer, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.remotePeerPools[peerAddr]; !ok {
		this.remotePeerPools[peerAddr] = newFunServantPeerPool(this.myIp,
			peerAddr, this.cf.PoolCapacity, this.cf.IdleTimeout)
		this.remotePeerPools[peerAddr].Open()
	}

	return this.remotePeerPools[peerAddr].Get()
}

// sticky request to remote peer servant by key
// return nil if I'm the servant for this key
func (this *Proxy) StickyServant(key string) (peer *FunServantPeer, peerAddr string) {
	peerAddr = this.selector.PickPeer(key)
	if peerAddr == "" {
		return
	}

	log.Debug("sticky key[%s] servant peer: %s", key, peerAddr)

	svt, _ := this.remotePeerPools[peerAddr].Get()
	return svt, peerAddr
}

// get all other servants in the cluster
// FIXME lock, but can't dead lock with this.Servant()
func (this *Proxy) ClusterServants() map[string]*FunServantPeer {
	rv := make(map[string]*FunServantPeer)
	for peerAddr, _ := range this.remotePeerPools {
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
	for addr, pool := range this.remotePeerPools {
		m[addr] = pool.pool.StatsJSON()
	}

	pretty, _ := json.MarshalIndent(m, "", "    ")
	return string(pretty)
}

func (this *Proxy) StatsMap() map[string]string {
	m := make(map[string]string)
	for addr, pool := range this.remotePeerPools {
		m[addr] = pool.pool.StatsJSON()
	}

	return m
}
