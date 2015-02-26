package redis

import (
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
