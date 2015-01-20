package servant

import (
	"github.com/funkygao/metrics"
	"sync/atomic"
)

type servantStats struct {
	calls         metrics.PercentCounter
	callsFromPeer int64
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
