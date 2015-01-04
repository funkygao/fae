package proxy

import (
	"hash/adler32"
)

type StandardPeerSelector struct {
	peerAddrs []string // array of peerAddr, self inclusive
}

func newStandardPeerSelector() *StandardPeerSelector {
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
