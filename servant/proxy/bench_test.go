package proxy

import (
	"errors"
	"testing"
)

func BenchmarkServantPing(b *testing.B) {
	proxy := NewWithDefaultConfig()
	servant, _ := proxy.ServantByAddr(":9001")
	defer servant.Transport.Close()

	ctx := servant.NewContext("bench_ping", nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.Ping(ctx)
	}

	b.SetBytes(10)
}

func BenchmarkServantMcSet(b *testing.B) {
	proxy := NewWithDefaultConfig()
	servant, _ := proxy.ServantByAddr(":9001")
	defer servant.Transport.Close()

	ctx := servant.NewContext("bench_mcset", nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.McGet(ctx, "pool", "key")
	}
	b.SetBytes(10)
}

func BenchmarkIsIoError(b *testing.B) {
	b.ReportAllocs()
	err := errors.New("broken pipe")
	for i := 0; i < b.N; i++ {
		IsIoError(err)
	}
}

func BenchmarkIsNotIoError(b *testing.B) {
	b.ReportAllocs()
	err := errors.New("blah")
	for i := 0; i < b.N; i++ {
		IsIoError(err)
	}
}
