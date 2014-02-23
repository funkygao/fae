package proxy

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

// A kind of pool.Resource
type FunServantPeer struct {
	*rpc.FunServantClient

	pool *funServantPeerPool
}

func (this *FunServantPeer) Close() {
	this.Transport.Close()
}

func (this *FunServantPeer) Recycle() {
	if this.Transport.IsOpen() {
		this.pool.Put(this)
	} else {
		this.pool.Put(nil)
	}
}
