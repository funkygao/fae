package peer

import (
	"encoding/json"
	"github.com/funkygao/golib/ip"
	log "github.com/funkygao/log4go"
	"net"
	"sync"
	"time"
)

type peerMessage map[string]interface{}

func (this *peerMessage) marshal() (data []byte, err error) {
	data, err = json.Marshal(this)
	return
}

func (this *peerMessage) unmarshal(data []byte) (err error) {
	err = json.Unmarshal(data, this)
	return
}

type Peer struct {
	rwmutex *sync.RWMutex

	c      net.Conn
	picker PeerPicker

	selfAddr          string
	groupAddr         string
	gaddr             *net.UDPAddr
	heartbeatInterval int
	deadThreshold     float64

	neighbors map[string]time.Time
}

func NewPeer(gaddr string, interval int,
	deadThreshold float64, replicas int) (this *Peer) {
	this = new(Peer)
	this.rwmutex = new(sync.RWMutex)
	this.groupAddr = gaddr
	this.selfAddr = ip.LocalIpv4Addrs()[0]
	this.picker = newPeerPicker(this.selfAddr, replicas)
	this.heartbeatInterval = interval
	this.deadThreshold = deadThreshold
	this.neighbors = make(map[string]time.Time)
	return
}

func (this *Peer) Neighbors() *map[string]time.Time {
	this.rwmutex.RLock()
	defer this.rwmutex.RUnlock()
	return &this.neighbors
}

func (this *Peer) killNeighbor(ip string) {
	this.rwmutex.Lock()
	defer this.rwmutex.Unlock()

	delete(this.neighbors, ip)
	this.picker.DelPeer(ip)
	log.Info("Peer[%s] killed", ip)

	log.Debug("Neighbors: %+v", this.neighbors)
}

func (this *Peer) refreshNeighbor(ip string) {
	this.rwmutex.Lock()
	defer this.rwmutex.Unlock()

	if _, present := this.neighbors[ip]; !present {
		log.Info("Peer[%s] joined", ip)
		this.picker.AddPeer(ip)
	}

	this.neighbors[ip] = time.Now()

	log.Debug("Neighbors: %+v", this.neighbors)
}

func (this *Peer) PickServer(key string) (serverAddr string, ok bool) {
	return this.picker.PickPeer(key)
}

func (this *Peer) Start() (err error) {
	this.gaddr, err = net.ResolveUDPAddr("udp4", this.groupAddr)
	if err != nil {
		return
	}
	this.c, err = net.ListenMulticastUDP("udp4", nil, this.gaddr)
	if err != nil {
		return
	}

	go this.runHeartbeat()
	go this.discoverPeers()

	log.Info("Self[%s] joined at %s", this.selfAddr, this.groupAddr)

	return
}

func (this *Peer) publish(msg peerMessage) (err error) {
	var body []byte
	body, err = msg.marshal()
	if err != nil {
		return
	}

	_, err = this.c.(*net.UDPConn).WriteToUDP(append(body, '\n'), this.gaddr)
	return
}
