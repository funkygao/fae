// stress test of fae
// simulate 3000 concurrent php-fpm requests for several rounds
package main

import (
	"flag"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"log"
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

	Basic = CallPing | CallIdGen
)

var (
	report stats
	ctx    *rpc.Context

	Concurrency     int
	Rounds          int // sessions=Rounds*Concurrency
	LoopsPerSession int // calls=sessions*LoopsPerSession
	Cmd             int
	host            string
	verbose         int
)

func init() {
	ctx = rpc.NewContext()
	ctx.Caller = "POST+/facebook/getPaymentRequestId/+34ca2cf6"
	ctx.Host = thrift.StringPtr("stress.test.local")
	ctx.Ip = thrift.StringPtr("127.0.0.1")
	ctx.Sid = thrift.StringPtr("bcf8f619")

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func parseFlag() {
	flag.IntVar(&LoopsPerSession, "loop", 1, "loops for each session")
	flag.IntVar(&Concurrency, "c", 3000, "concurrent num")
	flag.IntVar(&Cmd, "x", CallPing, "bitwise rpc calls")
	flag.IntVar(&Rounds, "n", 10, "rounds")
	flag.StringVar(&host, "h", "localhost", "rpc server host")
	flag.IntVar(&verbose, "v", 0, "verbose level")
	flag.Parse()
}

func main() {
	parseFlag()

	proxy := proxy.New(Concurrency, time.Minute*60)
	tryServantPool(proxy)

	time.Sleep(time.Second * 2)

	go report.run()

	wg := new(sync.WaitGroup)
	t1 := time.Now()
	for i := 0; i < Rounds; i++ {
		log.Printf("Round %3d", i)

		for j := 0; j < Concurrency; j++ {
			wg.Add(1)
			go runSession(proxy, wg, i, j)
		}

		wg.Wait()
	}

	elapsed := time.Since(t1)
	sessions := Rounds * Concurrency
	calls := sessions * LoopsPerSession
	log.Printf("Elapsed: %s, calls: %.1f/s, sessions: %.1f/s",
		elapsed,
		float64(calls)/elapsed.Seconds(),
		float64(sessions)/elapsed.Seconds())
}
