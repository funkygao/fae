package mysql

import (
	"fmt"
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
)

type StandardServerSelector struct {
	conf    *config.ConfigMysql
	clients map[string]*mysql // key is pool
}

func newStandardServerSelector(cf *config.ConfigMysql) (this *StandardServerSelector) {
	this = new(StandardServerSelector)
	this.conf = cf
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
	const SHARD_BASE_NUM = 200000 // TODO move the config
	bucket := fmt.Sprintf("%s%d", pool, (hintId/SHARD_BASE_NUM)+1)
	my, present := this.clients[bucket]
	if !present {
		return nil, ErrServerNotFound
	}

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
