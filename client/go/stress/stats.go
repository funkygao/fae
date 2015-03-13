package main

import (
	"github.com/funkygao/golib/gofmt"
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

type stats struct {
	concurrentN int32
	sessionN    int32 // aggregated sessions
	callErrs    int64
	callOk      int64
	connErrs    int64
	ioErrs      int64
}

func (this *stats) incCallErr() {
	atomic.AddInt64(&this.callErrs, 1)
}

func (this *stats) incCallOk() {
	atomic.AddInt64(&this.callOk, 1)
}

func (this *stats) incSessions() {
	atomic.AddInt32(&this.sessionN, 1)
}

func (this *stats) incConnErrs() {
	atomic.AddInt64(&this.connErrs, 1)
}

func (this *stats) incIoErrs() {
	atomic.AddInt64(&this.ioErrs, 1)
}

func (this *stats) updateConcurrency(delta int32) {
	atomic.AddInt32(&this.concurrentN, delta)
}

func (this *stats) run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastCalls int64
	for _ = range ticker.C {
		if neatStat {
			log.Printf("c:%6d qps:%20s errs:%10s",
				Concurrency,
				gofmt.Comma(this.callOk-lastCalls),
				gofmt.Comma(this.callErrs))
		} else {
			log.Printf("c:%d sessions:%d conns:%d calls:%s qps:%s errs:%s go:%d",
				Concurrency,
				atomic.LoadInt32(&this.sessionN),
				atomic.LoadInt32(&this.concurrentN),
				gofmt.Comma(atomic.LoadInt64(&this.callOk)),
				gofmt.Comma(this.callOk-lastCalls),
				gofmt.Comma(this.callErrs),
				runtime.NumGoroutine())
		}

		lastCalls = this.callOk
	}

}
