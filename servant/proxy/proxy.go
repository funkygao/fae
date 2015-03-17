package proxy

import (
	"encoding/json"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/ip"
	log "github.com/funkygao/log4go"
	"sync"
	"time"
)

// Proxy of remote servant so that we can dispatch request
// to cluster instead of having to serve all by ourselves.
type Proxy struct {
	mutex sync.Mutex
	cf    *config.ConfigProxy
	myIp  string

	clusterTopologyReady bool
	clusterTopologyChan  chan bool

	testOnBorrow    func()                         // TODO
	remotePeerPools map[string]*funServantPeerPool // key is peerAddr, self not inclusive
	selector        PeerSelector
}

func New(cf *config.ConfigProxy) *Proxy {
	ips, _ := ip.LocalIpv4Addrs()
	this := &Proxy{
		cf:                   cf,
		remotePeerPools:      make(map[string]*funServantPeerPool),
		selector:             newStandardPeerSelector(),
		myIp:                 ips[0],
		clusterTopologyReady: false,
		clusterTopologyChan:  make(chan bool),
	}

	return this
}

func NewWithDefaultConfig() *Proxy {
	return New(config.NewDefaultProxy())
}

func NewWithPoolCapacity(capacity int) *Proxy {
	cf := config.NewDefaultProxy()
	cf.PoolCapacity = capacity
	return New(cf)
}

// auto fills Rid field.
func NewContext(reason string) *rpc.Context {
	ctx := rpc.NewContext()
	ctx.Reason = reason
	ctx.Rid = time.Now().UnixNano() // TODO strict enough?
	return ctx
}

func (this *Proxy) Enabled() bool {
	return this.cf.Enabled()
}

func (this *Proxy) StartMonitorCluster() {
	if !this.Enabled() {
		log.Warn("servant proxy disabled by proxy config section")
		return
	}

	peersChan := make(chan []string, 10)
	go etclib.WatchService(etclib.SERVICE_FAE, peersChan)

	for {
		select {
		case <-peersChan:
			peers, err := etclib.ServiceEndpoints(etclib.SERVICE_FAE)
			if err == nil {
				if !this.clusterTopologyReady {
					this.clusterTopologyReady = true
					close(this.clusterTopologyChan)
				}

				if len(peers) == 0 {
					// TODO panic?
					log.Warn("Empty cluster fae peers")
				} else {
					// no lock, because running within 1 goroutine
					this.selector.SetPeersAddr(peers...)
					this.refreshPeers(peers)

					log.Info("Cluster latest fae nodes: %+v", peers)
				}
			} else {
				log.Error("Cluster peers: %s", err)
			}
		}
	}

	// should never get here
	log.Warn("Cluster peers monitor died")
}

func (this *Proxy) AwaitClusterTopologyReady() {
	if this.clusterTopologyReady {
		return
	}

	<-this.clusterTopologyChan
}

func (this *Proxy) refreshPeers(peers []string) {
	// add all latest peers
	for _, peerAddr := range peers {
		this.addRemotePeerIfNecessary(peerAddr)
	}

	// kill died peers
	for peerAddr, _ := range this.remotePeerPools {
		alive := false
		for _, p := range peers {
			if p == peerAddr {
				// still alive
				alive = true
				break
			}
		}

		if !alive {
			log.Warn("peer[%s] gone away", peerAddr)

			this.mutex.Lock()
			this.remotePeerPools[peerAddr].Close() // kill all conns in this pool
			delete(this.remotePeerPools, peerAddr)
			this.mutex.Unlock()
		}
	}
}

func (this *Proxy) addRemotePeerIfNecessary(peerAddr string) {
	if peerAddr == this.cf.SelfAddr {
		return
	}

	this.mutex.Lock()

	if _, present := this.remotePeerPools[peerAddr]; !present {
		this.remotePeerPools[peerAddr] = newFunServantPeerPool(this.myIp,
			peerAddr, *this.cf)
		this.remotePeerPools[peerAddr].Open()
	}

	this.mutex.Unlock()
}

// Get or create a fae peer servant based on peer address
func (this *Proxy) ServantByAddr(peerAddr string) (svt *FunServantPeer, err error) {
	this.addRemotePeerIfNecessary(peerAddr)

	this.remotePeerPools[peerAddr].nextTxn()
	for i := 0; i < this.cf.PoolCapacity; i++ {
		svt, err = this.remotePeerPools[peerAddr].Get()
		if err == nil {
			break
		} else {
			log.Error("ServantByAddr[%s].%d: %v", peerAddr, i, err)

			if svt != nil {
				if IsIoError(err) {
					svt.Close()
				}
				svt.Recycle()
			}
		}
	}

	return
}

// sticky request to remote peer servant by key
// return nil if I'm the servant for this key
func (this *Proxy) ServantByKey(key string) (svt *FunServantPeer, err error) {
	peerAddr := this.selector.PickPeer(key)
	if peerAddr == this.cf.SelfAddr {
		return nil, nil
	}

	this.remotePeerPools[peerAddr].nextTxn()
	for i := 0; i < this.cf.PoolCapacity; i++ {
		svt, err = this.remotePeerPools[peerAddr].Get()
		if err == nil {
			// found a valid conn
			break
		} else {
			log.Error("ServantByKey[%s].%d: %v", key, i, err)

			if svt != nil {
				if IsIoError(err) {
					svt.Close()
				}
				svt.Recycle()
			}
		}
	}

	return
}

// Simulate a simple load balance
func (this *Proxy) RandServant() (svt *FunServantPeer, err error) {
	peerAddr := this.selector.RandPeer()
	if peerAddr == this.cf.SelfAddr {
		return nil, nil
	}

	this.remotePeerPools[peerAddr].nextTxn()
	for i := 0; i < this.cf.PoolCapacity; i++ {
		svt, err = this.remotePeerPools[peerAddr].Get()
		if err == nil {
			break
		} else {
			log.Error("RandServant.%d: %v", i, err)

			if svt != nil {
				if IsIoError(err) {
					svt.Close()
				}
				svt.Recycle()
			}

		}
	}

	return
}

// Remote only, self not inclusive
func (this *Proxy) RemoteServants(haltOnErr bool) ([]*FunServantPeer, error) {
	r := make([]*FunServantPeer, 0)
	for addr, pool := range this.remotePeerPools {
		pool.nextTxn()

		svt, err := pool.Get()
		if err != nil {
			if svt != nil {
				if IsIoError(err) {
					svt.Close()
				}
				svt.Recycle()
			}

			if haltOnErr {
				return nil, err
			} else {
				log.Error("RemoteServants[%s]: %s", addr, err)
			}
		} else {
			r = append(r, svt)
		}
	}

	return r, nil
}

// peer addresses in the cluster
func (this *Proxy) ClusterPeers() []string {
	addrs := make([]string, 0, len(this.remotePeerPools))
	for addr, _ := range this.remotePeerPools {
		addrs = append(addrs, addr)
	}
	return addrs
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

func (this *Proxy) Warmup() {
	this.AwaitClusterTopologyReady()

	t0 := time.Now()
	for _, peerPool := range this.remotePeerPools {
		for i := 0; i < this.cf.PoolCapacity; i++ {
			svt, err := peerPool.Get()
			if svt != nil {
				if err != nil {
					svt.Close()
				}
				svt.Recycle()
			}
		}
	}

	log.Debug("Proxy warmup within %s", time.Since(t0))
}
