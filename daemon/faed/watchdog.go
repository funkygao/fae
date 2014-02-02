package main

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/golib/gofmt"
	"runtime"
	"time"
)

func runWatchdog(ticker *time.Ticker) {
	startTime := time.Now()
	ms := new(runtime.MemStats)

	for _ = range ticker.C {
		runtime.ReadMemStats(ms)

		log.Info("ver:%s, tick:%ds goroutine:%d, mem:%s, elapsed:%s",
			BuildID,
			options.tick,
			runtime.NumGoroutine(),
			gofmt.ByteSize(ms.Alloc),
			time.Since(startTime))
	}
}
