package proxy

import (
	"hash/adler32"
	"math/rand"
	"time"
)

type StandardPeerSelector struct {
	peerAddrs []string // array of peerAddr, self inclusive
}

func newStandardPeerSelector() *StandardPeerSelector {
	rand.Seed(time.Now().UnixNano())
	return &StandardPeerSelector{}
}

func (this *StandardPeerSelector) SetPeersAddr(peerAddrs ...string) {
	this.peerAddrs = peerAddrs
}

func (this *StandardPeerSelector) PickPeer(key string) (peerAddr string) {
	// adler32 is almost same as crc32, but much 3 times faster
	checksum := adler32.Checksum([]byte(key))
	index := int(checksum) % len(this.peerAddrs)

	return this.peerAddrs[index]
}

func (this *StandardPeerSelector) RandPeer() string {
	return this.peerAddrs[rand.Perm(len(this.peerAddrs))[0]]
}
