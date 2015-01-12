package game

import (
	"github.com/funkygao/assert"
	"math"
	"math/rand"
	"testing"
)

func TestName3(t *testing.T) {
	slots := 3
	nm := newNameGen(slots)
	v1 := nm.Next()
	v2 := nm.Next()
	assert.NotEqual(t, v1, v2)
	t.Logf("v1:%s, v2:%s, slots:%d, mem usage: %d", v1, v2, slots, nm.Size())

	for i := 0; i < 10; i++ {
		t.Logf("%s", nm.Next())
	}
	t.Logf("%+v", nm.bits)
}

func TestNotDuplicatedName(t *testing.T) {
	nm := newNameGen(3)
	total := 0

	for i := 0; i < int(math.Pow(float64(NameCharMax-NameCharMin), 3)); i++ {
		v := nm.Next()
		//t.Logf("%s", v)
		assert.Equal(t, true, nm.Contains(v), i)

		total++
	}

	t.Logf("total: %d, bits: %+v", total, nm.bits)
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
