package servant

import (
	"github.com/funkygao/metrics"
)

type servantStats struct {
	calls metrics.PercentCounter
}

func (this *servantStats) registerMetrics() {
	this.calls = metrics.NewPercentCounter()
	metrics.Register("servant.calls", this.calls)
}

func (this *servantStats) inc(key string) {
	this.calls.Inc(key, 1)
}
