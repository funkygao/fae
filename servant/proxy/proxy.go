package proxy

import (
	"encoding/json"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/ip"
	log "github.com/funkygao/log4go"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
)

// Proxy of remote servant so that we can dispatch request
// to cluster instead of having to serve all by ourselves.
type Proxy struct {
	mutex sync.Mutex
	cf    config.ConfigProxy
	myIp  string

	remotePeerPools map[string]*funServantPeerPool // key is peerAddr
	selector        PeerSelector
}

func New(cf config.ConfigProxy) *Proxy {
	this := &Proxy{
		cf:              cf,
		remotePeerPools: make(map[string]*funServantPeerPool),
		selector:        newStandardPeerSelector(),
		myIp:            ip.LocalIpv4Addrs()[0],
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

	var zkConnected = true
	for zkConnected {
		select {
		case <-peersChan:
			peers, err := etclib.ServiceEndpoints(etclib.SERVICE_FAE)
			if err == nil {
				// no lock, because running within 1 goroutine
				log.Trace("Cluster latest fae nodes: %+v", peers)

				this.selector.SetPeersAddr(peers)
				this.refreshPeers(peers)
			} else {
				log.Error("Cluster peers: %s", err)

				if err == zk.ErrClosing || err == zk.ErrConnectionClosed {
					zkConnected = false
				}
			}
		}
	}

	log.Warn("Cluster peers monitor died")
}

func (this *Proxy) refreshPeers(peers []string) {
	// add all latest peers
	for _, peerAddr := range peers {
		if peerAddr == this.cf.SelfAddr {
			continue
		}

		this.addRemotePeerIfNecessary(peerAddr)
	}

	// kill died peers
	for peerAddr, _ := range this.remotePeerPools {
		if peerAddr == this.cf.SelfAddr {
			continue
		}

		alive := false
		for _, p := range peers {
			if p == peerAddr {
				// still alive
				alive = true
				break
			}
		}

		if !alive {
			log.Trace("peer[%s] gone away", peerAddr)

			this.mutex.Lock()
			this.remotePeerPools[peerAddr].Close() // kill all conns in this pool
			delete(this.remotePeerPools, peerAddr)
			this.mutex.Unlock()
		}
	}

}

func (this *Proxy) addRemotePeerIfNecessary(peerAddr string) {
	this.mutex.Lock()

	if _, present := this.remotePeerPools[peerAddr]; !present {
		this.remotePeerPools[peerAddr] = newFunServantPeerPool(this.myIp,
			peerAddr, this.cf)
		this.remotePeerPools[peerAddr].Open()
	}

	this.mutex.Unlock()
}

// Get or create a fae peer servant based on peer address
func (this *Proxy) Servant(peerAddr string) (*FunServantPeer, error) {
	this.addRemotePeerIfNecessary(peerAddr)

	log.Debug("servant by addr[%s]: {txn: %d}", peerAddr,
		this.remotePeerPools[peerAddr].nextTxn())
	return this.remotePeerPools[peerAddr].Get()
}

// sticky request to remote peer servant by key
// return nil if I'm the servant for this key
func (this *Proxy) ServantByKey(key string) (*FunServantPeer, error) {
	peerAddr := this.selector.PickPeer(key)
	if peerAddr == this.cf.SelfAddr {
		return nil, nil
	}

	log.Debug("sevant by key[%s]: {peer: %s, txn: %d}", key, peerAddr,
		this.remotePeerPools[peerAddr].nextTxn())
	return this.remotePeerPools[peerAddr].Get()
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
