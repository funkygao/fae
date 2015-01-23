// +build !plan9,!windows

package servant

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/couch"
	"github.com/funkygao/fae/servant/game"
	"github.com/funkygao/fae/servant/memcache"
	"github.com/funkygao/fae/servant/mongo"
	"github.com/funkygao/fae/servant/mysql"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/fae/servant/redis"
	"github.com/funkygao/fae/servant/store"
	"github.com/funkygao/golib/cache"
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/golib/idgen"
	"github.com/funkygao/golib/mutexmap"
	"github.com/funkygao/golib/server"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"labix.org/v2/mgo"
	"net/http"
	"reflect"
	"regexp"
	"sync/atomic"
	"time"
)

type FunServantImpl struct {
	conf *config.ConfigServant

	digitNormalizer  *regexp.Regexp
	phpReasonPercent metrics.PercentCounter // user's behavior

	// stateful mem data related to services
	mysqlMergeMutexMap *mutexmap.MutexMap
	dbCacheStore       store.Store
	dbCacheHits        metrics.PercentCounter

	sessionN int64           // total sessions served since boot
	sessions *cache.LruCache // state kept for sessions FIXME kill it
	stats    *servantStats   // stats

	proxy *proxy.Proxy // remote fae agent

	idgen *idgen.IdGenerator   // global id generator
	game  *game.Game           // game engine
	lc    *cache.LruCache      // local cache
	mc    *memcache.ClientPool // memcache pool, auto sharding by key
	mg    *mongo.Client        // mongodb pool, auto sharding by shardId
	my    *mysql.MysqlCluster  // mysql pool, auto sharding by shardId
	rd    *redis.Client        // redis pool, auto sharding by pool name
	cb    *couch.Client        // couchbase client
}

func NewFunServant(cf *config.ConfigServant) (this *FunServantImpl) {
	this = &FunServantImpl{conf: cf}
	this.digitNormalizer = regexp.MustCompile(`\d+`)

	// http REST to export internal state
	if server.Launched() {
		server.RegisterHttpApi("/s/{cmd}",
			func(w http.ResponseWriter, req *http.Request,
				params map[string]interface{}) (interface{}, error) {
				return this.handleHttpQuery(w, req, params)
			}).Methods("GET")
	}

	this.sessions = cache.NewLruCache(cf.SessionMaxItems)
	this.mysqlMergeMutexMap = mutexmap.New(cf.Mysql.JsonMergeMaxOutstandingItems)

	// stats
	this.stats = new(servantStats)
	this.stats.registerMetrics()
	this.phpReasonPercent = metrics.NewPercentCounter()
	metrics.Register("php.reason", this.phpReasonPercent)
	this.dbCacheHits = metrics.NewPercentCounter()
	metrics.Register("db.cache.hits", this.dbCacheHits)

	// proxy can dynamically auto discover peers
	if this.conf.Proxy.Enabled() {
		log.Debug("creating servant: proxy")
		this.proxy = proxy.New(this.conf.Proxy)
	} else {
		panic("peers proxy required")
	}

	this.createServants()

	return
}

func (this *FunServantImpl) Start() {
	go this.showStats()
	go this.proxy.StartMonitorCluster()
	go this.watchConfigReloaded()

	this.warmUp()
}

func (this *FunServantImpl) Flush() {
	log.Debug("servants flushing...")
	// TODO
	this.my.Close()
	log.Trace("servants flushed")
}

func (this *FunServantImpl) AddErr(n int64) {
	this.stats.addErr(n)
}

func (this *FunServantImpl) createServants() {
	log.Info("creating servants...")

	log.Debug("creating servant: idgen")
	var err error
	this.idgen, err = idgen.NewIdGenerator(this.conf.IdgenWorkerId)
	if err != nil {
		panic(err)
	}

	log.Debug("creating servant: game")
	this.game = game.New(this.conf.Game)

	if this.conf.Lcache.Enabled() {
		log.Debug("creating servant: lcache")
		this.lc = cache.NewLruCache(this.conf.Lcache.MaxItems)
		this.lc.OnEvicted = this.onLcLruEvicted
	}

	if this.conf.Memcache.Enabled() {
		log.Debug("creating servant: memcache")
		this.mc = memcache.New(this.conf.Memcache)
	}

	if this.conf.Redis.Enabled() {
		log.Debug("creating servant: redis")
		this.rd = redis.New(this.conf.Redis)
	}

	if this.conf.Mysql.Enabled() {
		log.Debug("creating servant: mysql")
		this.my = mysql.New(this.conf.Mysql)
	}

	switch this.conf.Mysql.CacheStore {
	case "mem":
		this.dbCacheStore = store.NewMemStore(this.conf.Mysql.CacheStoreMemMaxItems)

	case "redis":
		this.dbCacheStore = store.NewRedisStore(this.conf.Mysql.CacheStoreRedisPool,
			this.conf.Redis)

	default:
		panic("unknown cache store")
	}

	if this.conf.Mongodb.Enabled() {
		log.Debug("creating servant: mongodb")
		this.mg = mongo.New(this.conf.Mongodb)
		if this.conf.Mongodb.DebugProtocol ||
			this.conf.Mongodb.DebugHeartbeat {
			mgo.SetLogger(&mongoProtocolLogger{})
			mgo.SetDebug(this.conf.Mongodb.DebugProtocol)
		}
	}

	if this.conf.Couchbase.Enabled() {
		log.Debug("creating servant: couchbase")

		var err error
		// pool is always 'default'
		this.cb, err = couch.New(this.conf.Couchbase.Servers, "default")
		if err != nil {
			log.Error("couchbase: %s", err)
		}
	}

	log.Info("servants created")
}

// TODO kill some servant if new conf turns it off
func (this *FunServantImpl) recreateServants(cf *config.ConfigServant) {
	log.Info("recreating servants...")

	if this.conf.IdgenWorkerId != cf.IdgenWorkerId {
		log.Debug("recreating servant: idgen")
		var err error
		this.idgen, err = idgen.NewIdGenerator(cf.IdgenWorkerId)
		if err != nil {
			panic(err)
		}
	}

	if !reflect.DeepEqual(*this.conf.Game, *cf.Game) {
		log.Debug("recreating servant: game")
		this.game = game.New(cf.Game)
	}

	if cf.Lcache.Enabled() &&
		this.conf.Lcache.MaxItems != cf.Lcache.MaxItems {
		log.Debug("recreating servant: lcache")
		this.lc = cache.NewLruCache(cf.Lcache.MaxItems)
		this.lc.OnEvicted = this.onLcLruEvicted
	}

	if cf.Memcache.Enabled() &&
		!reflect.DeepEqual(*this.conf.Memcache, *cf.Memcache) {
		log.Debug("recreating servant: memcache")
		this.mc = memcache.New(cf.Memcache)
	}

	if cf.Redis.Enabled() &&
		!reflect.DeepEqual(*this.conf.Redis, *cf.Redis) {
		log.Debug("recreating servant: redis")
		this.rd = redis.New(cf.Redis)
	}

	if cf.Mysql.Enabled() &&
		!reflect.DeepEqual(*this.conf.Mysql, *cf.Mysql) {
		log.Debug("recreating servant: mysql")
		this.my = mysql.New(cf.Mysql)

		switch cf.Mysql.CacheStore {
		case "mem":
			this.dbCacheStore = store.NewMemStore(cf.Mysql.CacheStoreMemMaxItems)

		case "redis":
			this.dbCacheStore = store.NewRedisStore(cf.Mysql.CacheStoreRedisPool,
				cf.Redis)

		default:
			panic("unknown cache store")
		}
	}

	if cf.Mongodb.Enabled() &&
		!reflect.DeepEqual(*this.conf.Mongodb, *cf.Mongodb) {
		log.Debug("recreating servant: mongodb")
		this.mg = mongo.New(cf.Mongodb)
	}

	if cf.Couchbase.Enabled() &&
		!reflect.DeepEqual(*this.conf.Couchbase, *cf.Couchbase) {
		log.Debug("recreating servant: couchbase")

		var err error
		// pool is always 'default'
		this.cb, err = couch.New(cf.Couchbase.Servers, "default")
		if err != nil {
			log.Error("couchbase: %s", err)
		}
	}

	this.conf = cf

	log.Info("servants recreated")
}

func (this *FunServantImpl) Runtime() map[string]interface{} {
	r := make(map[string]interface{})
	r["sessions"] = atomic.LoadInt64(&this.sessionN)
	r["call.errs"] = this.stats.callsErr
	r["call.peer.from"] = this.stats.callsFromPeer
	r["call.peer.to"] = this.stats.callsToPeer
	for _, key := range this.stats.calls.Keys() {
		r["call["+key+"]"] = this.stats.calls.Percent(key)
	}
	for _, key := range this.phpReasonPercent.Keys() {
		r["reason["+key+"]"] = this.phpReasonPercent.Percent(key)
	}
	for _, key := range this.dbCacheHits.Keys() {
		r["dbcache["+key+"]"] = this.dbCacheHits.Percent(key)
	}

	return r
}

func (this *FunServantImpl) showStats() {
	ticker := time.NewTicker(config.Engine.Servants.StatsOutputInterval)
	defer ticker.Stop()

	// TODO show most recent stats, reset at some interval

	for _ = range ticker.C {
		callsN := this.stats.calls.Total()
		log.Info("rpc: {sessions:%s, calls:%s, avg:%.1f; errs:%s peer.from:%s, peer.to:%s}",
			gofmt.Comma(this.sessionN),
			gofmt.Comma(callsN),
			float64(callsN)/float64(this.sessionN+1), // +1 to avoid divide by zero
			gofmt.Comma(this.stats.callsErr),
			gofmt.Comma(this.stats.callsFromPeer),
			gofmt.Comma(this.stats.callsToPeer))
	}
}

func (this *FunServantImpl) watchConfigReloaded() {
	for {
		select {
		case cf := <-config.Engine.ReloadedChan:
			this.recreateServants(cf.Servants)
		}
	}
}
