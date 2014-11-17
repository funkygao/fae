package couch

import (
	couchbase "github.com/couchbaselabs/go-couchbase"
	"sync"
)

type Client struct {
	mutex   sync.Mutex
	pool    couchbase.Pool
	buckets map[string]*couchbase.Bucket
}

func New(endpoint string, pool string) (this *Client, err error) {
	c, e := couchbase.Connect(endpoint)
	if e != nil {
		err = e
		return
	}

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
