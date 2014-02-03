package engine

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

	c net.Conn

	localAddr         string
	groupAddr         string
	gaddr             *net.UDPAddr
	heartbeatInterval int

	neighbors map[string]bool
}

func newPeer(gaddr string, interval int) (this *Peer) {
	this = new(Peer)
	this.RWMutex = new(sync.RWMutex)
	this.groupAddr = gaddr
	this.heartbeatInterval = interval
	this.neighbors = make(map[string]bool)
	return
}

func (this *Peer) markNeighbor(ip string, alive bool) {
	this.Lock()
	defer this.Unlock()

	this.neighbors[ip] = alive
}

func (this *Peer) Start() (err error) {
	this.localAddr = ip.LocalIpv4Addrs()[0]
	this.gaddr, _ = net.ResolveUDPAddr("udp4", this.groupAddr)
	this.c, err = net.ListenMulticastUDP("udp4", nil, this.gaddr)
	if err != nil {
		return
	}

	go this.runHeartbeat()
	go this.recvMessages()

	log.Info("Peer[%s] joined at %s", this.localAddr, this.groupAddr)

	return
}

func (this *Peer) recvMessages() {
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
			log.Error(err)
			continue
		}

		this.handleMessage(msg)
	}
}

func (this *Peer) handleMessage(msg peerMessage) {
	neighborIp, present := msg["ip"]
	if !present {
		return
	}

	this.markNeighbor(neighborIp.(string), true)
}

func (this *Peer) Publish(msg peerMessage) (err error) {
	var body []byte
	body, err = msg.marshal()
	if err != nil {
		return
	}

	log.Debug("publish %+v", msg)

	_, err = this.c.(*net.UDPConn).WriteToUDP(append(body, '\n'), this.gaddr)
	return
}

func (this *Peer) runHeartbeat() {
	t := time.NewTicker(time.Duration(this.heartbeatInterval) * time.Second)
	defer t.Stop()

	var msg = peerMessage{}
	var err error
	msg["ip"] = this.localAddr
	for _ = range t.C {
		if err = this.Publish(msg); err != nil {
			log.Error(err)
		}
	}
}
