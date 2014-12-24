package store

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/server"
	"testing"
)

func TestMemStore(t *testing.T) {
	s := NewMemStore(100)
	runStoreTest(t, s)
}

func TestRedisStore(t *testing.T) {
	svr := server.NewServer("test")
	svr.LoadConfig("../../etc/faed.cf.sample")
	section, _ := svr.Conf.Section("servants.redis")
	cf := &config.ConfigRedis{}
	cf.LoadConfig(section)

	s := NewRedisStore(cf)
	runStoreTest(t, s)
}

func runStoreTest(t *testing.T, s Store) {
	key := "hello"
	val, present := s.Get(key)
	assert.Equal(t, false, present)
	assert.Equal(t, nil, val)

	s.Set(key, "world")
	val, present = s.Get(key)
	assert.Equal(t, true, present)
	assert.Equal(t, "world", val)

	s.Del(key)
	val, present = s.Get(key)
	assert.Equal(t, false, present)
	assert.Equal(t, nil, val)
}

func BenchmarkRedisStoreSet(b *testing.B) {
	b.ReportAllocs()

	svr := server.NewServer("test")
	svr.LoadConfig("../../etc/faed.cf.sample")
	section, _ := svr.Conf.Section("servants.redis")
	cf := &config.ConfigRedis{}
	cf.LoadConfig(section)

	s := NewRedisStore(cf)
	k, v := "hello_benchmark", "world_benchmark"
	for i := 0; i < b.N; i++ {
		s.Set(k, v)
	}

	b.SetBytes(int64(len(k + v)))
}

func BenchmarkRedisStoreGet(b *testing.B) {
	b.ReportAllocs()

	svr := server.NewServer("test")
	svr.LoadConfig("../../etc/faed.cf.sample")
	section, _ := svr.Conf.Section("servants.redis")
	cf := &config.ConfigRedis{}
	cf.LoadConfig(section)

	s := NewRedisStore(cf)
	k, v := "hello_benchmark", "world_benchmark"
	s.Set(k, v)
	for i := 0; i < b.N; i++ {
		s.Get(k)
	}

	b.SetBytes(int64(len(k + v)))
}
