package redis

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/server"
	"github.com/funkygao/msgpack"
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
	s.LoadConfig("../../etc/faed.cf.sample")
	section, _ := s.Conf.Section("servants.redis")
	cf := &config.ConfigRedis{}
	cf.LoadConfig(section)

	var (
		pool = "default"
		val  interface{}
		err  error
	)

	c := New(cf)
	val, err = c.Get(pool, "hello")
	assert.Equal(t, ErrorDataNotExists.Error(), err.Error())

	err = c.Set(pool, "hello", "world")
	assert.Equal(t, nil, err)
	val, err = c.Get(pool, "hello")
	assert.Equal(t, nil, err)
	assert.Equal(t, "world", string(val.([]byte)))

	err = c.Del(pool, "hello")
	assert.Equal(t, nil, err)
	val, err = c.Get(pool, "hello")
	assert.Equal(t, ErrorDataNotExists, err)

	err = c.Del(pool, "hello") // del again
	assert.Equal(t, nil, err)

	var fooValue interface{}
	_, err = c.Call("SET", pool, "foo", "bar")
	assert.Equal(t, nil, err)
	fooValue, err = c.Call("GET", pool, "foo")
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", string(fooValue.([]byte)))
}

type customType struct {
	X    int // must be public when using msgpack for serialization
	Y    int
	Name string
}

func TestComplexDataType(t *testing.T) {
	s := server.NewServer("test")
	s.LoadConfig("../../etc/faed.cf.sample")
	section, _ := s.Conf.Section("servants.redis")
	cf := &config.ConfigRedis{}
	cf.LoadConfig(section)

	var (
		pool = "default"
		//val  interface{}
		err error
	)

	c := New(cf)
	name := "funky.gao"
	data := customType{X: 12, Y: 87, Name: name}
	encodedData, err := msgpack.Marshal(data)
	assert.Equal(t, nil, err)
	key := "custome_complex_data_type"
	err = c.Set(pool, key, encodedData)
	assert.Equal(t, nil, err)
	val, err := c.Get(pool, key)
	assert.Equal(t, nil, err)
	var val1 customType
	msgpack.Unmarshal(val.([]byte), &val1)
	assert.Equal(t, name, val1.Name)

}
