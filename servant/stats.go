package servant

import (
	"github.com/funkygao/metrics"
	"log"
	"os"
	"time"
)

type servantStats struct {
	TotalFailedCalls metrics.Counter
}

func (this *servantStats) Start(interval time.Duration) {
	if interval > 0 {
		go metrics.Log(metrics.DefaultRegistry,
			interval, log.New(os.Stderr, "", log.LstdFlags))
	}
}
