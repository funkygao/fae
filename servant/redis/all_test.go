package redis

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/server"
	"net"
	"testing"
)

// 536 ns/op TODO
func BenchmarkResolveTCPAddr(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		net.ResolveTCPAddr("tcp", "12.2.11.1:6378")
	}
}

func TestCRUD(t *testing.T) {
	s := server.NewServer("test")
	s.LoadConfig("../../etc/faed.cf")
	section, _ := s.Conf.Section("servants.redis")
	cf := &config.ConfigRedis{}
	cf.LoadConfig(section)

	var (
		pool = "default"
		val  string
		err  error
	)

	c := New(cf)
	err = c.Get(pool, "hello", &val)
	assert.Equal(t, ErrorDataNotExists.Error(), err.Error())

	err = c.Set(pool, "hello", "world")
	assert.Equal(t, nil, err)
	err = c.Get(pool, "hello", &val)
	assert.Equal(t, nil, err)
	assert.Equal(t, "world", val)

	err = c.Del(pool, "hello")
	assert.Equal(t, nil, err)
	err = c.Get(pool, "hello", &val)
	assert.Equal(t, ErrorDataNotExists, err)

	err = c.Del(pool, "hello") // del again
	assert.Equal(t, nil, err)
}
