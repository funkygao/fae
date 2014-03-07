package main

import (
	"github.com/funkygao/golib/gofmt"
	"log"
	"sync/atomic"
	"time"
)

type stats struct {
	concurrentN int32
	sessionN    int32 // aggregated sessions
	totalCalls  int64
	callErrs    int64
	callOk      int64
	connErrs    int64
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

func (this *stats) modifyConcurrency(delta int32) {
	atomic.AddInt32(&this.concurrentN, delta)
}

func (this *stats) run() {
	var lastCalls int64
	for {
		if lastCalls != 0 {
			log.Printf("sessions: %5d concurrency: %5d calls:%12s cps: %9s errs: %9s",
				this.sessionN,
				this.concurrentN,
				gofmt.Comma(this.callOk),
				gofmt.Comma(this.callOk-lastCalls),
				gofmt.Comma(this.callErrs))
		} else {
			log.Printf("sessions: %5d concurrency: %5d calls: %12s",
				this.sessionN,
				this.concurrentN,
				gofmt.Comma(this.callOk))
		}

		lastCalls = this.callOk

		time.Sleep(time.Second)
	}
}
