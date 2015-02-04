package proxy

type PeerSelector interface {
	SetPeersAddr(peerAddrs ...string) // self inclusive
	PickPeer(key string) string       // return peer addr, self inclusive
	RandPeer() string
}
