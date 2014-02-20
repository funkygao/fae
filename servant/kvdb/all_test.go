package kvdb

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer("server")
	err := s.Open()
	assert.Equal(t, nil, err)
}

func TestServlet(t *testing.T) {
	s := NewServlet("test")
	defer s.Close()

	key := []byte("hello")
	value := []byte("world")

	// open
	err := s.Open()
	assert.Equal(t, nil, err)

	// put
	err = s.Put(key, value)
	assert.Equal(t, nil, err)

	// get
	val, err := s.Get(key)
	assert.Equal(t, value, val)

	// del
	err = s.Delete(key)
	assert.Equal(t, nil, err)

	// after del, get again
	val, err = s.Get(key)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(val))
}

func BenchmarkServletPut(b *testing.B) {
	s := NewServlet("test")
	s.Open()
	defer s.Close()

	key := []byte("hello")
	value := []byte("world")
	for i := 0; i < b.N; i++ {
		s.Put(key, value)
	}
	b.SetBytes(int64(len(key) + len(value)))
}

func BenchmarkServletGet(b *testing.B) {
	s := NewServlet("test")
	s.Open()
	defer s.Close()

	key := []byte("hello")
	value := []byte("world")
	s.Put(key, value)
	for i := 0; i < b.N; i++ {
		s.Get(key)
	}
	b.SetBytes(int64(len(key) + len(value)))
}
