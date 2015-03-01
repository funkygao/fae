package servant

import (
	"github.com/funkygao/metrics"
	"sync/atomic"
)

var (
	svtStats servantStats
)

type servantStats struct {
	calls metrics.PercentCounter

	callsFromPeer int64
	callsToPeer   int64

	callsSlow int64 // TODO mv to engine
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

func (this *servantStats) incCallSlow() {
	atomic.AddInt64(&this.callsSlow, 1)
}
