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
	t.Logf("v1:%s, v2:%s, slots:%d, mem usage: %d", v1, v2, slots, nm.Size())
	t.Logf("%+v", nm.bits)
	for i := 0; i < 1000; i++ {
		t.Logf("%s", nm.Next())
	}
	assert.NotEqual(t, v1, v2)

}

// 529 ns/op
func BenchmarkNextName(b *testing.B) {
	b.ReportAllocs()
	nm := New(3)
	for i := 0; i < b.N; i++ {
		nm.Next()
	}
}
