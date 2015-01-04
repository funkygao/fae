package proxy

import (
	"github.com/funkygao/golib/hash"
)

type ConsistentPeerSelector struct {
	peers *hash.Map // array of peerAddr, self inclusive
}

func newConsistentPeerSelector() *ConsistentPeerSelector {
	return &ConsistentPeerSelector{peers: hash.New(32, nil)} // TODO 32?
}

func (this *ConsistentPeerSelector) SetPeersAddr(peerAddrs ...string) {
	this.peers.Add(peerAddrs...)
}

func (this *ConsistentPeerSelector) PickPeer(key string) (peerAddr string) {
	return this.peers.Get(key)
}
