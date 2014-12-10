package proxy

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/ip"
	"github.com/funkygao/golib/pool"
	log "github.com/funkygao/log4go"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// A kind of pool.Resource
type FunServantPeer struct {
	pool.Resource
	*rpc.FunServantClient

	pool *funServantPeerPool

	once sync.Once

	// ctx related
	rid  int64
	myIp string // self ip addr, set only once
}

func (this *FunServantPeer) Close() {
	this.Transport.Close()
	this.Recycle()
}

func (this *FunServantPeer) Recycle() {
	if this.Transport.IsOpen() {
		this.pool.Put(this)
	} else {
		this.pool.Put(nil)
	}
}

func (this *FunServantPeer) Addr() string {
	return this.pool.serverAddr
}

func (this *FunServantPeer) NewContext(reason string, uid *int64) *rpc.Context {
	ctx := rpc.NewContext() // TODO pool
	atomic.AddInt64(&this.rid, 1)
	ctx.Rid = strconv.FormatInt(this.rid, 10)
	ctx.Reason = reason
	ctx.Uid = uid
	this.once.Do(func() {
		this.myIp = ip.LocalIpv4Addrs()[0]
	})
	ctx.Host = this.myIp

	return ctx
}

// append my request id and my host ip to ctx
func (this *FunServantPeer) HijackContext(ctx *rpc.Context) {
	atomic.AddInt64(&this.rid, 1)
	ctx.Rid = ctx.Rid + ":" + strconv.FormatInt(this.rid, 10)
	this.once.Do(func() {
		this.myIp = ip.LocalIpv4Addrs()[0]
	})
	ctx.Host = ctx.Host + ":" + this.myIp
}

// a conn pool to a fae endpoint
type funServantPeerPool struct {
	serverAddr string

	capacity    int
	idleTimeout time.Duration
	pool        *pool.ResourcePool
}

func newFunServantPeerPool(serverAddr string, capacity int,
	idleTimeout time.Duration) (this *funServantPeerPool) {
	this = &funServantPeerPool{idleTimeout: idleTimeout, capacity: capacity,
		serverAddr: serverAddr}
	return
}

func (this *funServantPeerPool) connect(serverAddr string) (*rpc.FunServantClient,
	error) {
	transportFactory := thrift.NewTBufferedTransportFactory(4 << 10) // TODO
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(serverAddr) // should never timeout
	if err != nil {
		return nil, err
	}

	useTransport := transportFactory.GetTransport(transport)
	client := rpc.NewFunServantClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		log.Error("conn peer[%s]: %s", serverAddr, err)

		return nil, err
	}

	log.Trace("peer[%s] connected", serverAddr)
	return client, nil
}

func (this *funServantPeerPool) Open() {
	factory := func() (pool.Resource, error) {
		client, err := this.connect(this.serverAddr)
		if err != nil {
			return nil, err
		}

		return &FunServantPeer{FunServantClient: client, pool: this}, nil
	}

	this.pool = pool.NewResourcePool(factory,
		this.capacity, this.capacity,
		this.idleTimeout)
}

func (this *funServantPeerPool) Get() (*FunServantPeer, error) {
	fun, err := this.pool.Get()
	if err != nil {
		return nil, err
	}

	return fun.(*FunServantPeer), nil
}

func (this *funServantPeerPool) Put(conn *FunServantPeer) {
	if !this.pool.IsClosed() {
		this.pool.Put(conn)
	}
}
