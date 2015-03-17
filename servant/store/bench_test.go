package store

import (
	"testing"
)

func BenchmarkRedisStoreSet(b *testing.B) {
	b.ReportAllocs()

	s := getRedisStore()
	k, v := "hello_benchmark", "world_benchmark"
	for i := 0; i < b.N; i++ {
		s.Set(k, v)
	}

	b.SetBytes(int64(len(k + v)))
}

func BenchmarkRedisStoreGet(b *testing.B) {
	b.ReportAllocs()

	s := getRedisStore()
	k, v := "hello_benchmark", "world_benchmark"
	s.Set(k, v)
	for i := 0; i < b.N; i++ {
		s.Get(k)
	}

	b.SetBytes(int64(len(k + v)))
}
