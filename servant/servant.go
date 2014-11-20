// +build !plan9,!windows

package servant

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/couch"
	"github.com/funkygao/fae/servant/memcache"
	"github.com/funkygao/fae/servant/mongo"
	"github.com/funkygao/fae/servant/mysql"
	"github.com/funkygao/fae/servant/peer"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/cache"
	"github.com/funkygao/golib/idgen"
	"github.com/funkygao/golib/server"
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo"
	"net/http"
	"time"
)

type FunServantImpl struct {
	conf *config.ConfigServant

	sessions *cache.LruCache // state kept for sessions FIXME kill it
	stats    *servantStats   // stats

	proxy *proxy.Proxy         // remote fae agent
	peer  *peer.Peer           // topology of cluster
	idgen *idgen.IdGenerator   // global id generator
	lc    *cache.LruCache      // local cache
	mc    *memcache.ClientPool // memcache pool, auto sharding by key
	mg    *mongo.Client        // mongodb pool, auto sharding by shardId
	my    *mysql.MysqlCluster  // mysql pool, auto sharding by shardId
	cb    *couch.Client        // couchbase client
}

func NewFunServant(cf *config.ConfigServant) (this *FunServantImpl) {
	this = &FunServantImpl{conf: cf}
	this.sessions = cache.NewLruCache(20 << 10) // TODO config

	// stats
	this.stats = new(servantStats)
	this.stats.registerMetrics()

	// http REST
	if server.Launched() {
		server.RegisterHttpApi("/s/{cmd}",
			func(w http.ResponseWriter, req *http.Request,
				params map[string]interface{}) (interface{}, error) {
				return this.handleHttpQuery(w, req, params)
			}).Methods("GET")
	}

	// remote fae peer
	if this.conf.PeerEnabled() {
		this.peer = peer.NewPeer(this.conf.PeerGroupAddr,
			this.conf.PeerHeartbeatInterval,
			this.conf.PeerDeadThreshold, this.conf.PeersReplica)
		this.proxy = proxy.New(this.conf.Proxy.PoolCapacity,
			this.conf.Proxy.IdleTimeout)
	}

	// idgen, always present
	this.idgen = idgen.NewIdGenerator(this.conf.DataCenterId, this.conf.AgentId)

	// local cache
	if this.conf.Lcache.Enabled() {
		this.lc = cache.NewLruCache(this.conf.Lcache.LruMaxItems)
		this.lc.OnEvicted = this.onLcLruEvicted
	}

	// memcache
	if this.conf.Memcache.Enabled() {
		this.mc = memcache.New(this.conf.Memcache)
	}

	// mysql
	if this.conf.Mysql.Enabled() {
		this.my = mysql.New(this.conf.Mysql)
	}

	// mongodb
	if this.conf.Mongodb.Enabled() {
		this.mg = mongo.New(this.conf.Mongodb)

		if this.conf.Mongodb.DebugProtocol ||
			this.conf.Mongodb.DebugHeartbeat {
			mgo.SetLogger(&mongoProtocolLogger{})
			mgo.SetDebug(this.conf.Mongodb.DebugProtocol)
		}
	}

	if this.conf.Couchbase != nil &&
		this.conf.Couchbase.Servers != nil &&
		len(this.conf.Couchbase.Servers) > 0 {
		var err error
		// pool is always 'default'
		this.cb, err = couch.New(this.conf.Couchbase.Servers, "default")
		if err != nil {
			log.Error("couchbase: %s", err)
		}
	}

	return
}

func (this *FunServantImpl) Start() {
	this.warmUp()
	go this.showStats()

	if this.peer != nil {
		if err := this.peer.Start(); err != nil {
			log.Error("peer start: %v", err)
			this.peer = nil
		}
	}

}

func (this *FunServantImpl) showStats() {
	ticker := time.NewTicker(config.Servants.StatsOutputInterval)
	defer ticker.Stop()

	for _ = range ticker.C {
		log.Info("rpc: {sessions:%d, calls:%d, avg:%d}",
			this.sessions.Len(),
			this.stats.calls.Total(),
			this.stats.calls.Total()/int64(this.sessions.Len()))
	}
}
