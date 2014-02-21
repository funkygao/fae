package servant

import (
	"github.com/funkygao/metrics"
)

type servantStats struct {
	calls    metrics.PercentCounter
	inBytes  metrics.Counter
	outBytes metrics.Counter
}

func (this *servantStats) registerMetrics() {
	this.calls = metrics.NewPercentCounter()
	metrics.Register("servant.calls", this.calls)
	this.inBytes = metrics.NewCounter()
	metrics.Register("servant.in.bytes", this.inBytes)
	this.outBytes = metrics.NewCounter()
	metrics.Register("servant.out.bytes", this.outBytes)
}

func (this *servantStats) inc(key string) {
	this.calls.Inc(key, 1)
}
