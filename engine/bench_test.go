package engine

import (
	"github.com/funkygao/thrift/lib/go/thrift"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkSelectChanDefault(b *testing.B) {
	quit := make(chan struct{})
	for i := 0; i < b.N; i++ {
		select {
		case <-quit:
		default:
		}
	}
}

func BenchmarkTimeNow(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkTimeSince(b *testing.B) {
	b.ReportAllocs()
	t1 := time.Now()
	for i := 0; i < b.N; i++ {
		time.Since(t1)
	}
}

func BenchmarkDispatcher(b *testing.B) {
	b.ReportAllocs()
	handler := func(client thrift.TTransport) {}
	dispatcher := newRpcDispatcher(true, 10000, handler)
	for i := 0; i < b.N; i++ {
		dispatcher.Dispatch(nil)
	}
}

func BenchmarkEngineStatsActiveSessionN(b *testing.B) {
	var activeSessionN int64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomic.AddInt64(&activeSessionN, 1)
	}
}

func BenchmarkEngineStatsCallPerSecond(b *testing.B) {
	s := newEngineStats()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.CallPerSecond.Mark(1)
	}
}

func BenchmarkEngineStatsCallLatencies(b *testing.B) {
	b.ReportAllocs()
	s := newEngineStats()
	elapsed := time.Second
	for i := 0; i < b.N; i++ {
		s.CallLatencies.Update(elapsed.Nanoseconds() / 1e6)
	}
}

func BenchmarkStatsRuntime(b *testing.B) {
	b.ReportAllocs()
	s := newEngineStats()
	for i := 0; i < b.N; i++ {
		s.Runtime()
	}
}
