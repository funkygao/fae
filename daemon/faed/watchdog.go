package main

import (
	"github.com/funkygao/fae/engine"
	"github.com/funkygao/golib/gofmt"
	log "github.com/funkygao/log4go"
	"runtime"
	"syscall"
	"time"
)

func runWatchdog(interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var (
		startTime    = time.Now()
		ms           = new(runtime.MemStats)
		rusage       = &syscall.Rusage{}
		lastUserTime int64
		lastSysTime  int64
		userTime     int64
		sysTime      int64
		userCpuUtil  float64
		sysCpuUtil   float64
	)

	for _ = range ticker.C {
		runtime.ReadMemStats(ms)

		syscall.Getrusage(syscall.RUSAGE_SELF, rusage)
		userTime = rusage.Utime.Sec*1000000000 + int64(rusage.Utime.Usec)
		sysTime = rusage.Stime.Sec*1000000000 + int64(rusage.Stime.Usec)
		userCpuUtil = float64(userTime-lastUserTime) * 100 / float64(interval)
		sysCpuUtil = float64(sysTime-lastSysTime) * 100 / float64(interval)

		lastUserTime = userTime
		lastSysTime = sysTime

		log.Info("ver:%s, tick:%ds goroutine:%d, mem:%s, elapsed:%s",
			engine.BuildID,
			options.tick,
			runtime.NumGoroutine(),
			gofmt.ByteSize(ms.Alloc),
			time.Since(startTime))
		log.Info("cpu: %3.2f%% us, %3.2f%% sy, rss:%s",
			userCpuUtil,
			sysCpuUtil,
			gofmt.ByteSize(float64(rusage.Maxrss*1024)))
	}
}
