package kvdb

import (
	"github.com/funkygao/assert"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer("server", 0)
	err := s.Open()
	defer os.RemoveAll("server")
	assert.Equal(t, nil, err)
}

func TestServlet(t *testing.T) {
	s := newServlet("test")
	defer s.close()
	defer os.RemoveAll("test")

	key := []byte("hello")
	value := []byte("world")

	// open
	err := s.open()
	assert.Equal(t, nil, err)

	// put
	err = s.put(key, value)
	assert.Equal(t, nil, err)

	// get
	val, err := s.get(key)
	assert.Equal(t, value, val)

	// del
	err = s.delete(key)
	assert.Equal(t, nil, err)

	// after del, get again
	val, err = s.get(key)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(val))
}

func BenchmarkServletPut(b *testing.B) {
	s := newServlet("test")
	s.open()
	defer s.close()
	defer os.RemoveAll("test")

	key := []byte("hello")
	value := []byte("world")
	for i := 0; i < b.N; i++ {
		s.put(key, value)
	}
	b.SetBytes(int64(len(key) + len(value)))
}

func BenchmarkServletGet(b *testing.B) {
	s := newServlet("test")
	s.open()
	defer s.close()
	defer os.RemoveAll("test")

	key := []byte("hello")
	value := []byte("world")
	s.put(key, value)
	for i := 0; i < b.N; i++ {
		s.get(key)
	}
	b.SetBytes(int64(len(key) + len(value)))
}

func BenchmarkServletDelete(b *testing.B) {
	b.ReportAllocs()
	s := newServlet("test")
	s.open()
	defer s.close()
	defer os.RemoveAll("test")

	key := []byte("hello")
	value := []byte("world")
	s.put(key, value)
	for i := 0; i < b.N; i++ {
		s.delete(key)
	}
}
