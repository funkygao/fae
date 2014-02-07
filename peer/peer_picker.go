package peer

import (
	"github.com/funkygao/golib/hash"
	"sync"
)

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	AddPeer(peers ...string)
	DelPeer(peers ...string)
	PickPeer(key string) (serverAddr string, ok bool)
}

type consistentPeerPicker struct {
	self string
	*sync.Mutex
	peers *hash.Map
}

func newPeerPicker(self string) PeerPicker {
	return &consistentPeerPicker{peers: hash.New(3, nil), Mutex: new(sync.Mutex),
		self: self}
}

func (this *consistentPeerPicker) DelPeer(peers ...string) {
	// TODO
}

func (this *consistentPeerPicker) AddPeer(peers ...string) {
	this.Lock()
	defer this.Unlock()
	this.peers.Add(peers...)
}

func (this *consistentPeerPicker) PickPeer(key string) (serverAddr string, ok bool) {
	this.Lock()
	defer this.Unlock()
	if this.peers.IsEmpty() {
		return "", false
	}

	if peer := this.peers.Get(key); peer != this.self {
		return peer, true
	}

	return "", false

}
