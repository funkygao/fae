package memcache

import (
	"github.com/funkygao/assert"
	"os/exec"
	"testing"
	"time"
)

func TestMemcacheClient(t *testing.T) {
	cmd := exec.Command("memcached")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Memcached start: %v", err)
	}
	defer cmd.Process.Kill()
	time.Sleep(time.Second)

	mc := newMemcacheClient()
	err := mc.Connect("localhost:11211")
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}

	stored, err := mc.Set("test", 0, 0, []byte("Hello 中关村"))
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if !stored {
		t.Errorf("want true, got %v", stored)
	}

	results, err := mc.Get("test")
	assert.Equal(t, "Hello 中关村", string(results[0].Value))
}

func BenchmarkMemcacheClientSet(b *testing.B) {
	cmd := exec.Command("memcached")
	if err := cmd.Start(); err != nil {
		b.Fatalf("Memcached start: %v", err)
	}
	defer cmd.Process.Kill()
	time.Sleep(time.Second)

	mc := newMemcacheClient()
	err := mc.Connect("localhost:11211")
	if err != nil {
		b.Fatalf("Connect: %v", err)
	}

	for i := 0; i < b.N; i++ {
		stored, err := mc.Set("test", 0, 0, []byte("Hello 中关村"))
		if err != nil {
			b.Fatalf("Set: %v", err)
		}
		if !stored {
			b.Errorf("want true, got %v", stored)
		}
	}

}

func BenchmarkHash(b *testing.B) {
	key := "user:23424"
	for i := 0; i < b.N; i++ {
		findServer(key)
	}
}
