package couch

import (
	couchbase "github.com/couchbaselabs/go-couchbase"
	"sync"
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
func New(endpoint string, pool string) (this *Client, err error) {
	// connect to couchbase cluster: any node in the cluster is ok
	// internally: GET /pools
	c, e := couchbase.Connect(endpoint)
	if e != nil {
		err = e
		return
	}

	// internally: GET /pools/default, then GET /pools/default/buckets
	// get the vBucketServerMap and nodes ip:port in cluster
	p, e := c.GetPool(pool)
	if e != nil {
		err = e
		return
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
