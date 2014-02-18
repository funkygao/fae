package servant

import (
	"github.com/funkygao/fae/config"
	rest "github.com/funkygao/fae/http"
	"github.com/funkygao/fae/servant/memcache"
	"github.com/funkygao/fae/servant/mongo"
	"github.com/funkygao/fae/servant/peer"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/cache"
	"github.com/funkygao/golib/idgen"
	"labix.org/v2/mgo"
	"net/http"
	"time"
)

type FunServantImpl struct {
	conf *config.ConfigServant

	proxy *proxy.Proxy // remote fae agent
	peer  *peer.Peer   // topology of cluster
	idgen *idgen.IdGenerator
	lc    *cache.LruCache
	mc    *memcache.Client
	mg    *mongo.Client
}

func NewFunServant(cf *config.ConfigServant) (this *FunServantImpl) {
	this = &FunServantImpl{conf: cf}
	this.idgen = idgen.NewIdGenerator()

	if this.conf.Lcache.Enabled() {
		this.lc = cache.NewLruCache(this.conf.Lcache.LruMaxItems)
		this.lc.OnEvicted = this.onLcLruEvicted
	}

	if this.conf.Memcache.Enabled() {
		memcacheServers := this.conf.Memcache.ServerList()
		this.mc = memcache.New(this.conf.Memcache.HashStrategy, memcacheServers...)
		this.mc.Timeout = time.Duration(this.conf.Memcache.Timeout) * time.Second
		this.mc.MaxIdleConnsPerServer = this.conf.Memcache.MaxIdleConnsPerServer
	}

	if this.conf.Mongodb.Enabled() {
		this.mg = mongo.New(this.conf.Mongodb)
		if this.conf.Mongodb.DebugProtocol ||
			this.conf.Mongodb.DebugHeartbeat {
			mgo.SetLogger(&mongoProtocolLogger{})
			mgo.SetDebug(this.conf.Mongodb.DebugProtocol)
		}
	}

	if this.conf.PeerEnabled() {
		this.peer = peer.NewPeer(this.conf.PeerGroupAddr,
			this.conf.PeerHeartbeatInterval,
			this.conf.PeerDeadThreshold, this.conf.PeersReplica)
		this.proxy = proxy.New()
	}

	if rest.Launched() {
		rest.RegisterHttpApi("/s/{cmd}",
			func(w http.ResponseWriter, req *http.Request,
				params map[string]interface{}) (interface{}, error) {
				return this.handleHttpQuery(w, req, params)
			}).Methods("GET")
	}

	return
}

func (this *FunServantImpl) Start() {
	go this.runWatchdog()
	if this.peer != nil {
		this.peer.Start()
	}
}
