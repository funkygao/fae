package couchbasebench

import (
	"testing"

	couchbase "github.com/couchbaselabs/go-couchbase"
)

func mf(err error, msg string) {
	if err != nil {
		println(err)
	}
}

func BenchmarkSimpleSet(b *testing.B) {
	b.ReportAllocs()

	c, err := couchbase.Connect("http://localhost:8091/")

	p, err := c.GetPool("default")
	mf(err, "pool")

	bucket, err := p.GetBucket("default")
	mf(err, "bucket")
	for i := 0; i < b.N; i++ {
		bucket.Set(",k", 90, map[string]interface{}{"x": 1})
	}
}
