package engine

import (
	"fmt"
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/metrics"
	"log"
	"os"
	"runtime"
	"time"
)

type engineStats struct {
	startedAt time.Time
	memStats  *runtime.MemStats

	TotalFailedCalls    metrics.Counter
	TotalFailedSessions metrics.Counter
	TotalSlowSessions   metrics.Counter
	TotalSlowCalls      metrics.Counter
	SessionLatencies    metrics.Histogram
	CallLatencies       metrics.Histogram
	SessionPerSecond    metrics.Meter
	CallPerSecond       metrics.Meter
	CallPerSession      metrics.Histogram
}

func newEngineStats() (this *engineStats) {
	this = new(engineStats)
	this.memStats = new(runtime.MemStats)
	this.registerMetrics()
	return
}

func (this *engineStats) registerMetrics() {
	this.TotalFailedSessions = metrics.NewCounter()
	metrics.Register("total.sessions.fail", this.TotalFailedSessions)
	this.TotalSlowSessions = metrics.NewCounter()
	metrics.Register("total.sessions.slow", this.TotalSlowSessions)
	this.TotalFailedCalls = metrics.NewCounter()
	metrics.Register("total.calls.fail", this.TotalFailedCalls)
	this.TotalSlowCalls = metrics.NewCounter()
	metrics.Register("total.calls.slow", this.TotalSlowCalls)
	this.SessionLatencies = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.session", this.SessionLatencies)
	this.CallLatencies = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("latency.call", this.CallLatencies)
	this.SessionPerSecond = metrics.NewMeter()
	metrics.Register("rps.session", this.SessionPerSecond)
	this.CallPerSecond = metrics.NewMeter()
	metrics.Register("rps.call", this.CallPerSecond)
	this.CallPerSession = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("call.per.session", this.CallPerSession)
}

// TODO
func (this engineStats) String() string {
	return ""
}

func (this *engineStats) Start(t time.Time, interval time.Duration, logfile string) {
	this.startedAt = t

	metricsWriter, err := os.OpenFile(logfile,
		os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	if interval > 0 {
		metrics.Log(metrics.DefaultRegistry,
			interval, log.New(metricsWriter, "", log.LstdFlags))
	}
}

func (this *engineStats) Runtime() map[string]interface{} {
	runtime.ReadMemStats(this.memStats)

	s := make(map[string]interface{})
	s["goroutines"] = runtime.NumGoroutine()
	s["memory.allocated"] = gofmt.ByteSize(this.memStats.Alloc).String()
	s["memory.mallocs"] = gofmt.ByteSize(this.memStats.Mallocs).String()
	s["memory.frees"] = gofmt.ByteSize(this.memStats.Frees).String()
	s["memory.last_gc"] = this.memStats.LastGC
	s["memory.gc.num"] = this.memStats.NumGC
	s["memory.gc.num_per_second"] = float64(this.memStats.NumGC) / time.
		Since(this.startedAt).Seconds()
	s["memory.gc.total_pause"] = fmt.Sprintf("%dms",
		this.memStats.PauseTotalNs/uint64(time.Millisecond))
	s["memory.heap.alloc"] = gofmt.ByteSize(this.memStats.HeapAlloc).String()
	s["memory.heap.sys"] = gofmt.ByteSize(this.memStats.HeapSys).String()
	s["memory.heap.idle"] = gofmt.ByteSize(this.memStats.HeapIdle).String()
	s["memory.heap.released"] = gofmt.ByteSize(this.memStats.HeapReleased).String()
	s["memory.heap.objects"] = gofmt.Comma(int64(this.memStats.HeapObjects))
	s["memory.stack"] = gofmt.ByteSize(this.memStats.StackInuse).String()
	gcPausesMs := make([]string, 0, 20)
	for _, pauseNs := range this.memStats.PauseNs {
		if pauseNs == 0 {
			continue
		}

		pauseStr := fmt.Sprintf("%dms",
			pauseNs/uint64(time.Millisecond))
		if pauseStr == "0ms" {
			continue
		}

		gcPausesMs = append(gcPausesMs, pauseStr)
	}
	s["memory.gc.pauses"] = gcPausesMs

	return s
}
