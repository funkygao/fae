package main

import (
	"github.com/funkygao/fae/engine"
	"github.com/funkygao/golib/gofmt"
	log "github.com/funkygao/log4go"
	"runtime"
	"time"
)

func runWatchdog(ticker *time.Ticker) {
	startTime := time.Now()
	ms := new(runtime.MemStats)

	for _ = range ticker.C {
		runtime.ReadMemStats(ms)

		log.Info("ver:%s, tick:%ds goroutine:%d, mem:%s, elapsed:%s",
			engine.BuildID,
			options.tick,
			runtime.NumGoroutine(),
			gofmt.ByteSize(ms.Alloc),
			time.Since(startTime))
	}
}
