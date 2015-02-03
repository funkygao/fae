package proxy

import (
	"github.com/funkygao/golib/hash"
	"math/rand"
	"time"
)

type ConsistentPeerSelector struct {
	peerAddrs []string  // just for random selecting
	peers     *hash.Map // array of peerAddr, self inclusive
}

func newConsistentPeerSelector() *ConsistentPeerSelector {
	rand.Seed(time.Now().UnixNano())
	return &ConsistentPeerSelector{peers: hash.New(32, nil)} // TODO 32?
}

func (this *ConsistentPeerSelector) SetPeersAddr(peerAddrs ...string) {
	this.peerAddrs = peerAddrs
	this.peers.Add(peerAddrs...)
}

func (this *ConsistentPeerSelector) PickPeer(key string) (peerAddr string) {
	return this.peers.Get(key)
}

func (this *ConsistentPeerSelector) RandPeer() string {
	return this.peerAddrs[rand.Perm(len(this.peerAddrs))[0]]
}
