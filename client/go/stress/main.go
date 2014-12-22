// stress test of fae
// simulate 3000 concurrent php-fpm requests for several rounds
package main

import (
	"flag"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/gofmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	CallPing = 1 << iota
	CallIdGen
	CallLCache
	CallMemcache
	CallMongo
	CallKvdb

	CallPingIdgen   = CallPing | CallIdGen
	CallIdgenLcache = CallIdGen | CallLCache

	MC_POOL = "default"
)

var (
	report stats
	ctx    *rpc.Context

	SampleRate      int
	Concurrency     int
	Rounds          int // sessions=Rounds*Concurrency
	LoopsPerSession int // calls=sessions*LoopsPerSession
	Cmd             int
	host            string
	verbose         int
)

func init() {
	ctx = rpc.NewContext()
	ctx.Reason = "POST+/facebook/getPaymentRequestId/+34ca2cf6"
	ctx.Host = "stress.test.local"
	ctx.Ip = "127.0.0.1"
	ctx.Rid = "bcf8f619"

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func parseFlag() {
	flag.IntVar(&LoopsPerSession, "loop", 1, "loops for each session")
	flag.IntVar(&Concurrency, "c", 3000, "concurrent num")
	flag.IntVar(&SampleRate, "s", Concurrency, "sampling rate")
	flag.IntVar(&Cmd, "x", CallPing, "bitwise rpc calls")
	flag.IntVar(&Rounds, "n", 10, "rounds")
	flag.StringVar(&host, "h", "localhost", "rpc server host")
	flag.IntVar(&verbose, "v", 0, "verbose level")
	flag.Usage = showUsage
	flag.Parse()
}

func main() {
	parseFlag()

	cf := config.ConfigProxy{PoolCapacity: 20}
	proxy := proxy.New(cf)
	tryServantPool(proxy)

	time.Sleep(time.Second * 2)

	go report.run()

	wg := new(sync.WaitGroup)
	t1 := time.Now()
	for i := 0; i < Rounds; i++ {
		for j := 0; j < Concurrency; j++ {
			wg.Add(1)
			go runSession(proxy, wg, i+1, j)
		}

		wg.Wait()
	}

	elapsed := time.Since(t1)
	log.Printf("Elapsed: %s, calls: {%s, %.1f/s}, sessions: {%s, %.1f/s}, errors: {conn:%d, call:%d}",
		elapsed,
		gofmt.Comma(report.callOk),
		float64(report.callOk)/elapsed.Seconds(),
		gofmt.Comma(int64(report.sessionN)),
		float64(report.sessionN)/elapsed.Seconds(),
		report.connErrs,
		report.callErrs)
}
