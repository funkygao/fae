package mysql

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/cache"
	"github.com/funkygao/golib/str"
	log "github.com/funkygao/log4go"
	"strconv"
)

type StandardServerSelector struct {
	conf    *config.ConfigMysql
	clients map[string]*mysql // key is pool

	lookupCache *cache.LruCache // {pool+hintId: *mysql}
}

func newStandardServerSelector(cf *config.ConfigMysql) (this *StandardServerSelector) {
	this = new(StandardServerSelector)
	this.conf = cf
	this.lookupCache = cache.NewLruCache(cf.LookupCacheMaxItems)
	this.clients = make(map[string]*mysql)
	for _, server := range cf.Servers {
		my := newMysql(server.DSN(), &cf.Breaker)
		for retries := uint(0); retries < cf.Breaker.FailureAllowance; retries++ {
			log.Debug("mysql connecting: %s", server.DSN())

			if my.Open() == nil && my.Ping() == nil {
				// sql.Open() does not establish any connections to the database
				// sql.Ping() does
				break
			}

			my.breaker.Fail()
		}

		my.db.SetMaxIdleConns(cf.MaxIdleConnsPerServer)
		// https://code.google.com/p/go/source/detail?r=8a7ac002f840
		my.db.SetMaxOpenConns(cf.MaxConnsPerServer)
		this.clients[server.Pool] = my
	}

	return
}

func (this *StandardServerSelector) PickServer(pool string,
	table string, hintId int) (*mysql, error) {
	if this.shardedPool(pool) {
		return this.pickShardedServer(pool, table, hintId)
	}

	return this.pickNonShardedServer(pool, table)
}

func (this *StandardServerSelector) ServerByBucket(bucket string) (*mysql, error) {
	my, present := this.clients[bucket]
	if !present {
		return nil, ErrServerNotFound
	}

	return my, nil
}

func (this *StandardServerSelector) Servers() []*mysql {
	r := make([]*mysql, 0)
	for _, m := range this.clients {
		r = append(r, m)
	}

	return r
}

func (this *StandardServerSelector) shardedPool(pool string) bool {
	if _, present := this.conf.GlobalPools[pool]; present {
		return false
	}

	return true
}

func (this *StandardServerSelector) pickShardedServer(pool string,
	table string, hintId int) (*mysql, error) {
	const (
		sb1 = "SELECT shardId FROM "
		sb2 = " WHERE entityId=?"
	)

	// FIXME how to handle cache kick?
	// TODO itoa is too slow, 143 ns/op, use int as cache key
	key := pool + strconv.Itoa(hintId)
	if conn, present := this.lookupCache.Get(key); present {
		log.Debug("lookupCache[%s] hit", key)
		return conn.(*mysql), nil
	} else {
		log.Debug("lookupCache[%s] miss", key)
	}

	my, err := this.ServerByBucket(this.conf.LookupPool)
	if err != nil {
		return nil, err
	}

	sb := str.NewStringBuilder()
	sb.WriteString(sb1)
	sb.WriteString(this.conf.LookupTable(pool))
	sb.WriteString(sb2)
	rows, err := my.Query(sb.String(), hintId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// only 1 row in lookup table
	if !rows.Next() {
		return nil, ErrShardLookupNotFound
	}

	var shardId string
	if err = rows.Scan(&shardId); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	//bucket := fmt.Sprintf("%s%d", pool, (hintId/200000)+1)
	bucket := pool + shardId
	my, present := this.clients[bucket]
	if !present {
		return nil, ErrServerNotFound
	}

	//this.lookupCache.Set(key, my)

	return my, nil
}

func (this *StandardServerSelector) pickNonShardedServer(pool string,
	table string) (*mysql, error) {
	my, present := this.clients[pool]
	if !present {
		return nil, ErrServerNotFound
	}

	return my, nil
}

func (this *StandardServerSelector) endsWithDigit(pool string) bool {
	lastChar := pool[len(pool)-1]
	return lastChar >= '0' && lastChar <= '9'
}
