package couch

import (
	couchbase "github.com/couchbaselabs/go-couchbase"
	log "github.com/funkygao/log4go"
	"math/rand"
	"sync"
	"time"
)

// Couchbase is designed to be a drop-in replacement for an existing memcached server, while
// adding persistence, replication, failover and dynamic cluster reconfiguration.
type Client struct {
	mutex   sync.Mutex
	pool    couchbase.Pool
	buckets map[string]*couchbase.Bucket
}

// Till Couchbase 2.x releases, pool is a placeholder that doesn't have any special meaning
// Also note that no decisions have been made about what Couchbase will do with pools
func New(baseUrls []string, pool string) (this *Client, err error) {
	var (
		c couchbase.Client
		p couchbase.Pool
		e error
	)

	rand.Seed(time.Now().UTC().UnixNano())
	for _, i := range rand.Perm(len(baseUrls)) { // client side load balance
		// connect to couchbase cluster: any node in the cluster is ok
		// internally: GET /pools
		nodeUrl := baseUrls[i]
		c, e = couchbase.Connect(nodeUrl) // TODO timeout
		if e == nil {
			break
		}

		// failed to connect to this node in cluster
		log.Warn("couchbase[%s] connect fail: %s", nodeUrl, e.Error())
	}

	if e != nil {
		// max retry reached
		return nil, e
	}

	// internally: GET /pools/default, then GET /pools/default/buckets
	// get the vBucketServerMap and nodes ip:port in cluster
	// TODO connct to streamingUri, cluster updates are fetched from that conn
	p, e = c.GetPool(pool)
	if e != nil {
		return nil, e
	}

	this = new(Client)
	this.pool = p
	this.buckets = make(map[string]*couchbase.Bucket)
	return
}

// The unit of multi-tenancy in Couchbase is the “bucket”
// which represents a “virtual Couchbase Server instance” inside a single
// Couchbase Server cluster
// Bucket can be treated as database in mysql
// The limit of the number of buckets that can be configured within a cluster is 10
func (this *Client) GetBucket(bucket string) (*couchbase.Bucket, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	b, present := this.buckets[bucket]
	if present {
		return b, nil
	}

	b, e := this.pool.GetBucket(bucket)
	this.buckets[bucket] = b
	return b, e
}
