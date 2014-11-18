package mysql

import (
	"fmt"
	"github.com/funkygao/fae/config"
)

const (
	ServerActive      = "active"  // fully oprational
	ServerDead        = "dead"    // fully non-operational
	ServerPending     = "pending" // blocks clients, receives replicas
	ServerReplicating = "replica" // dead to clients, receive replicas, transfer vbuckets from one server to another
)

// vBucket-aware mysql cluster client selector.
//
// The vBucket mechanism provides a layer of indirection between the hashing algorithm
// and the server responsible for a given key.
//
// This indirection is useful in managing the orderly transition from one cluster configuration to
// another, whether the transition was planned (e.g. adding new servers to a cluster) or unexpected (e.g. a server failure)
//
// Every key belongs to a vBucket, which maps to a server instance
// The number of vBuckets is a fixed number
//
// key ---------------> server
// h(key) -> vBucket -> server
//
// servers = ['server1:11211', 'server2:11211', 'server3:11211']
// vbuckets = [0, 0, 1, 1, 2, 2]
// server_for_key(key) = servers[vbuckets[hash(key) % vbuckets.length]]
//
// how to add a new server:
// push config to all clients -> to make the new server useful, transfer vbuckets from one server to another, set
// them to ServerPending state on the receiving server
//
// The vBucket-Server map is updated internally: transmitted from server to all cluster participants:
// servers, clients and proxies
type VbucketServerSelector struct {
	conf    *config.ConfigMysql
	clients map[string]*mysql // key is pool
}

func newVbucketServerSelector(cf *config.ConfigMysql) (this *VbucketServerSelector) {
	this = new(VbucketServerSelector)
	this.conf = cf
	this.clients = make(map[string]*mysql)
	for _, server := range cf.Servers {
		my := newMysql(server.DSN(), &cf.Breaker)
		for retries := uint(0); retries < cf.Breaker.FailureAllowance; retries++ {
			if my.Open() == nil && my.Ping() == nil {
				// sql.Open() does not establish any connections to the database
				// it's lazy
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

func (this *VbucketServerSelector) PickServer(pool string,
	table string, hintId int) (*mysql, error) {
	if this.shardedPool(pool) {
		return this.pickShardedServer(pool, table, hintId)
	}

	return this.pickNonShardedServer(pool, table)
}

func (this *VbucketServerSelector) Servers() []*mysql {
	r := make([]*mysql, 0)
	for _, m := range this.clients {
		r = append(r, m)
	}

	return r

}

func (this *VbucketServerSelector) shardedPool(pool string) bool {
	switch pool {
	case "ShardLookup", "Global", "Tickets":
		return false
	default:
		return true
	}
}

func (this *VbucketServerSelector) pickShardedServer(pool string,
	table string, hintId int) (*mysql, error) {
	const SHARD_BASE_NUM = 200000 // TODO move the config
	bucket := fmt.Sprintf("%s%d", pool, (hintId/SHARD_BASE_NUM)+1)
	my, present := this.clients[bucket]
	if !present {
		return nil, ErrServerNotFound
	}

	return my, nil
}

func (this *VbucketServerSelector) pickNonShardedServer(pool string,
	table string) (*mysql, error) {
	my, present := this.clients[pool]
	if !present {
		return nil, ErrServerNotFound
	}

	return my, nil
}

func (this *VbucketServerSelector) endsWithDigit(pool string) bool {
	lastChar := pool[len(pool)-1]
	return lastChar >= '0' && lastChar <= '9'
}
