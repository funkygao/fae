package game

import (
	"github.com/funkygao/assert"
	"math/rand"
	"testing"
)

func TestPresence(t *testing.T) {
	p := newPresence()
	p.Update(12)
	p.Update(3)
	assert.Equal(t, []bool{true, true}, p.Onlines([]int64{12, 3}))
	assert.Equal(t, []bool{true, false}, p.Onlines([]int64{12, 31}))
	assert.Equal(t, []bool{false, false}, p.Onlines([]int64{122, 23}))
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
