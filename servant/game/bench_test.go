package game

import (
	"github.com/funkygao/fae/config"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkLockBasic(b *testing.B) {
	cf := &config.ConfigGame{
		LockMaxItems: 10,
		LockExpires:  10 * time.Second,
	}
	l := newLock(cf)
	k := "haha"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Lock(k)
		l.Unlock(k)
	}
}

// 827 ns/op
func BenchmarkNextName(b *testing.B) {
	b.ReportAllocs()
	nm := newNameGen(3)
	for i := 0; i < b.N; i++ {
		nm.Next()
	}
}

func BenchmarkRandInt31(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rand.Int31n(int32(NameCharMax - NameCharMin))
	}
}

// 40ms
func BenchmarkLoopPresence(b *testing.B) {
	m := make(map[int64]int)
	for i := 0; i < (1 << 20); i++ {
		m[rand.Int63()] = 0
	}

	for i := 0; i < b.N; i++ {
		for _, _ = range m {

		}
	}

}
