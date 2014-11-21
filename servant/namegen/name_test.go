package namegen

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestName3(t *testing.T) {
	slots := 3
	nm := New(slots)
	v1 := nm.Next()
	v2 := nm.Next()
	assert.NotEqual(t, v1, v2)
	t.Logf("v1:%s, v2:%s, slots:%d, mem usage: %d", v1, v2, slots, nm.Size())

	for i := 0; i < 10000; i++ {
		t.Logf("%s", nm.Next())
	}
	t.Logf("%+v", nm.bits)
}

func TestNotDuplicatedName(t *testing.T) {
	nm := New(3)
	total := 0
	for i := 0; i < int(NameCharMax-NameCharMin)*int(NameCharMax-NameCharMin)*int(NameCharMax-NameCharMin); i++ {
		v := nm.Next()
		//t.Logf("%s", v)
		assert.Equal(t, true, nm.Contains(v), i)

		total++
	}

	t.Logf("total: %d", total)
}

// 827 ns/op
func BenchmarkNextName(b *testing.B) {
	b.ReportAllocs()
	nm := New(3)
	for i := 0; i < b.N; i++ {
		nm.Next()
	}
}
