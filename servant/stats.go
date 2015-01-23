package servant

import (
	"github.com/funkygao/metrics"
	"sync/atomic"
)

type servantStats struct {
	calls         metrics.PercentCounter
	callsFromPeer int64
	callsToPeer   int64
	callsErr      int64
}

func (this *servantStats) registerMetrics() {
	this.calls = metrics.NewPercentCounter()
	metrics.Register("servant.calls", this.calls)
}

func (this *servantStats) inc(key string) {
	this.calls.Inc(key, 1)
}

func (this *servantStats) incPeerCall() {
	atomic.AddInt64(&this.callsFromPeer, 1)
}

func (this *servantStats) incCallPeer() {
	atomic.AddInt64(&this.callsToPeer, 1)
}

func (this *servantStats) incErr() {
	atomic.AddInt64(&this.callsErr, 1)
}

func (this *servantStats) addErr(n int64) {
	atomic.AddInt64(&this.callsErr, n)
}
