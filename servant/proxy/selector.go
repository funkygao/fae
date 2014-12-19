package proxy

type PeerSelector interface {
	SetPeersAddr(peerAddrs []string)
	PickPeer(key string) string // return peer addr
}
