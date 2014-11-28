package proxy

import (
	"testing"
)

func BenchmarkServantPing(b *testing.B) {
	proxy := New(10, 0)
	servant, _ := proxy.Servant(":9001")
	defer servant.Transport.Close()

	ctx := servant.NewContext("bench_ping", nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.Ping(ctx)

	}
	b.SetBytes(10)
}

func BenchmarkServantMcSet(b *testing.B) {
	proxy := New(10, 0)
	servant, _ := proxy.Servant(":9001")
	defer servant.Transport.Close()

	ctx := servant.NewContext("bench_mcset", nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.McGet(ctx, "pool", "key")
	}
	b.SetBytes(10)
}
