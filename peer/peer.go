package peer

import (
	"bufio"
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/funkygao/golib/ip"
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
	*sync.RWMutex

	c      net.Conn
	picker PeerPicker

	selfAddr          string
	groupAddr         string
	gaddr             *net.UDPAddr
	heartbeatInterval int
	deadThreshold     float64

	neighbors map[string]time.Time
}

func NewPeer(gaddr string, interval int, deadThreshold float64) (this *Peer) {
	this = new(Peer)
	this.RWMutex = new(sync.RWMutex)
	this.groupAddr = gaddr
	this.selfAddr = ip.LocalIpv4Addrs()[0]
	this.picker = newPeerPicker(this.selfAddr)
	this.heartbeatInterval = interval
	this.deadThreshold = deadThreshold
	this.neighbors = make(map[string]time.Time)
	return
}

func (this *Peer) Neighbors() *map[string]time.Time {
	this.RLock()
	defer this.RUnlock()
	return &this.neighbors
}

func (this *Peer) killNeighbor(ip string) {
	this.Lock()
	defer this.Unlock()

	delete(this.neighbors, ip)
	this.picker.DelPeer(ip)
	log.Info("Peer[%s] killed", ip)

	log.Debug("Neighbors: %+v", this.neighbors)
}

func (this *Peer) refreshNeighbor(ip string) {
	this.Lock()
	defer this.Unlock()

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

	log.Info("Peer[%s] joined at %s", this.selfAddr, this.groupAddr)

	return
}

func (this *Peer) discoverPeers() {
	defer func() {
		this.c.Close() // leave the multicast group
	}()

	var msg peerMessage
	reader := bufio.NewReader(this.c)
	for {
		// net.ListenMulticastUDP sets IP_MULTICAST_LOOP=0 as
		// default, so you never receive your own sent data
		// if you run both sender and receiver on (logically) same IP host
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Error(err)
			continue
		}

		if err := msg.unmarshal(line); err != nil {
			// Not our protocol, it may be SSDP or else
			continue
		}

		log.Debug("received peer: %+v", msg)

		neighborIp, present := msg["ip"]
		if !present {
			log.Info("Peer msg has no 'ip'")
			continue
		}

		this.refreshNeighbor(neighborIp.(string))
	}
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

func (this *Peer) runHeartbeat() {
	t := time.NewTicker(time.Duration(this.heartbeatInterval) * time.Second)
	defer t.Stop()

	var msg = peerMessage{}
	var err error
	msg["ip"] = this.selfAddr
	for _ = range t.C {
		if err = this.publish(msg); err != nil {
			log.Error("Publish fails: %v", err)
		}

		for ip, lastAccess := range this.neighbors {
			if time.Since(lastAccess).Seconds() > this.deadThreshold {
				// he is dead
				this.killNeighbor(ip)
			}
		}

	}
}
