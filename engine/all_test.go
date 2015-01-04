package engine

import (
	"testing"
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
