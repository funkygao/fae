package proxy

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/pool"
	"strconv"
)

// A single rpc client connection with remote peer
type FunServantPeer struct {
	pool.Resource
	*rpc.FunServantClient

	pool *funServantPeerPool
}

func newFunServantPeer(p *funServantPeerPool, c *rpc.FunServantClient) *FunServantPeer {
	this := new(FunServantPeer)
	this.FunServantClient = c
	this.pool = p
	return this
}

// TODO if peer conn broken, call this
func (this *FunServantPeer) Close() {
	this.Transport.Close()
	this.Recycle()
}

func (this *FunServantPeer) Recycle() {
	if this.Transport.IsOpen() {
		this.pool.pool.Put(this)
	} else {
		this.pool.pool.Put(nil)
	}
}

func (this *FunServantPeer) NewContext(reason string, uid *int64) *rpc.Context {
	ctx := rpc.NewContext() // TODO pool
	ctx.Rid = strconv.FormatInt(this.pool.nextTxn(), 10)
	ctx.Reason = reason
	ctx.Uid = uid
	ctx.Host = this.pool.myIp

	return ctx
}

// append my transaction id and my host ip to ctx
func (this *FunServantPeer) HijackContext(ctx *rpc.Context) {
	ctx.Rid = ctx.Rid + ":" + strconv.FormatInt(this.pool.nextTxn(), 10)
	ctx.Host = ctx.Host + ":" + this.pool.myIp
	ctx.Sticky = new(bool)
	*ctx.Sticky = true // tells peer it's from fae
}

func (this *FunServantPeer) Addr() string {
	return this.pool.peerAddr // peers in the pool share the remote peer addr
}
