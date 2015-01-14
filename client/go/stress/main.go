// stress test of fae
// simulate 3000 concurrent php-fpm requests for several rounds
package main

import (
	"flag"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/golib/server"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	ctx = rpc.NewContext()
	ctx.Reason = "stress.go"
	ctx.Host = "stress.test.local"
	ctx.Ip = "127.0.0.1"
	ctx.Rid = "bcf8f619"

	parseFlag()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	server.SetupLogging("var/test.log", "info", "var/panic.dump")
}

func parseFlag() {
	flag.IntVar(&LoopsPerSession, "loop", 1, "loops for each session")
	flag.IntVar(&Concurrency, "c", 3000, "concurrent num")
	flag.IntVar(&SampleRate, "s", Concurrency, "sampling rate")
	flag.IntVar(&Cmd, "x", CallDefault, "bitwise rpc calls")
	flag.IntVar(&Rounds, "n", 10, "rounds")
	flag.StringVar(&host, "host", "localhost", "rpc server host")
	flag.IntVar(&verbose, "v", 0, "verbose level")
	flag.StringVar(&zk, "zk", "localhost:2181", "zk server addr")
	flag.Usage = showUsage
	flag.Parse()
}

func main() {
	proxy := proxy.NewWithDefaultConfig()
	etclib.Dial([]string{zk})
	go proxy.StartMonitorCluster()
	proxy.AwaitClusterTopologyReady()

	// test pool
	testServantPool(proxy)
	pause("pool tested")

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
