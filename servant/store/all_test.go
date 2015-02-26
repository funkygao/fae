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
	s := getRedisStore()
	runStoreTest(t, s)
}

func getRedisStore() *RedisStore {
	svr := server.NewServer("test")
	svr.LoadConfig("../../etc/faed.cf.sample")
	section, _ := svr.Conf.Section("servants.redis")
	cf := &config.ConfigRedis{}
	cf.LoadConfig(section)

	return NewRedisStore("default", cf)
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
