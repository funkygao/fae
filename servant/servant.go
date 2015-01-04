// +build !plan9,!windows

package servant

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/couch"
	"github.com/funkygao/fae/servant/lock"
	"github.com/funkygao/fae/servant/memcache"
	"github.com/funkygao/fae/servant/mongo"
	"github.com/funkygao/fae/servant/mysql"
	"github.com/funkygao/fae/servant/namegen"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/fae/servant/redis"
	"github.com/funkygao/fae/servant/store"
	"github.com/funkygao/golib/cache"
	"github.com/funkygao/golib/idgen"
	"github.com/funkygao/golib/mutexmap"
	"github.com/funkygao/golib/server"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"labix.org/v2/mgo"
	"net/http"
	"regexp"
	"time"
)

type FunServantImpl struct {
	conf *config.ConfigServant

	digitNormalizer *regexp.Regexp

	// stateful mem data related to services
	mysqlMergeMutexMap *mutexmap.MutexMap
	dbCacheStore       store.Store
	dbCacheHits        metrics.PercentCounter

	sessionN int64           // total sessions served since boot
	sessions *cache.LruCache // state kept for sessions FIXME kill it
	stats    *servantStats   // stats

	// php client related
	phpLatency       metrics.Histogram      // in ms
	phpPayloadSize   metrics.Histogram      // in bytes
	phpReasonPercent metrics.PercentCounter // user's behavior

	// service drivers
	proxy   *proxy.Proxy         // remote fae agent
	idgen   *idgen.IdGenerator   // global id generator
	namegen *namegen.NameGen     // name generator
	lc      *cache.LruCache      // local cache
	mc      *memcache.ClientPool // memcache pool, auto sharding by key
	mg      *mongo.Client        // mongodb pool, auto sharding by shardId
	my      *mysql.MysqlCluster  // mysql pool, auto sharding by shardId
	rd      *redis.Client        // redis pool, auto sharding by pool name
	cb      *couch.Client        // couchbase client
	lk      *lock.Lock           // lock map
}

func NewFunServant(cf *config.ConfigServant) (this *FunServantImpl) {
	log.Debug("creating servants...")

	this = &FunServantImpl{conf: cf}
	this.sessions = cache.NewLruCache(cf.SessionEntries)

	this.mysqlMergeMutexMap = mutexmap.New(8 << 20) // 8M TODO
	this.digitNormalizer = regexp.MustCompile(`\d+`)

	// stats
	this.stats = new(servantStats)
	this.stats.registerMetrics()

	// record php latency histogram
	this.phpLatency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("php.latency", this.phpLatency)
	this.phpPayloadSize = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("php.payload", this.phpPayloadSize)
	this.phpReasonPercent = metrics.NewPercentCounter()
	metrics.Register("php.reason", this.phpReasonPercent)

	// http REST to export internal state
	if server.Launched() {
		server.RegisterHttpApi("/s/{cmd}",
			func(w http.ResponseWriter, req *http.Request,
				params map[string]interface{}) (interface{}, error) {
				return this.handleHttpQuery(w, req, params)
			}).Methods("GET")
	}

	log.Debug("creating servant: idgen")
	this.idgen = idgen.NewIdGenerator(this.conf.DataCenterId, this.conf.AgentId)

	log.Debug("creating servant: namegen")
	this.namegen = namegen.New(3)

	if this.conf.Proxy.Enabled() {
		log.Debug("creating servant: proxy")
		this.proxy = proxy.New(this.conf.Proxy)
	} else {
		panic("peers proxy disabled")
	}

	if this.conf.Lcache.Enabled() {
		log.Debug("creating servant: lcache")
		this.lc = cache.NewLruCache(this.conf.Lcache.MaxItems)
		this.lc.OnEvicted = this.onLcLruEvicted
	}

	if this.conf.Lock.Enabled() {
		log.Debug("creating servant: lock")
		this.lk = lock.New(this.conf.Lock)
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

	this.dbCacheHits = metrics.NewPercentCounter()
	metrics.Register("db.cache.hits", this.dbCacheHits)
	if this.conf.Mysql.CacheStore == "mem" {
		this.dbCacheStore = store.NewMemStore(this.conf.Mysql.CacheStoreMemMaxItems)
	} else if this.conf.Mysql.CacheStore == "redis" {
		this.dbCacheStore = store.NewRedisStore(this.conf.Mysql.CacheStoreRedisPool,
			this.conf.Redis)
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

	log.Debug("servants created")
	return
}

func (this *FunServantImpl) Start() {
	go this.showStats()
	go this.proxy.StartMonitorCluster()

	this.warmUp()
}

func (this *FunServantImpl) Flush() {
	log.Debug("servants flushing...")
	// TODO
	log.Trace("servants flushed")
}

func (this *FunServantImpl) showStats() {
	ticker := time.NewTicker(config.Engine.Servants.StatsOutputInterval)
	defer ticker.Stop()

	// TODO show most recent stats, reset at some interval

	for _ = range ticker.C {
		log.Info("rpc: {sessions:%d, calls:%d, avg:%d}",
			this.sessionN,
			this.stats.calls.Total(),
			this.stats.calls.Total()/int64(this.sessionN+1)) // +1 to avoid divide by zero
	}
}
