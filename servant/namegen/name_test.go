package namegen

import (
	"testing"
)

func TestAll(t *testing.T) {
	nm := New(3)
	t.Logf("%s", nm.Next())
}

func BenchmarkNextName(b *testing.B) {
	b.ReportAllocs()
	nm := New(4)
	b.Logf("mem space usage: %d", nm.Size())
	for i := 0; i < b.N; i++ {
		b.Logf("%s", nm.Next())
	}
}
