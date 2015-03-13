// stress test of fae
// simulate 3000 concurrent php-fpm requests for several rounds
package main

import (
	"flag"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/gofmt"
	"github.com/funkygao/golib/server"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	parseFlag()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	server.SetupLogging("stdout", "info", "", "", "")
}

func parseFlag() {
	flag.IntVar(&LoopsPerSession, "loop", 1, "loops for each session")
	flag.IntVar(&c1, "c", 1, "concurrency lower limit")
	flag.IntVar(&c2, "c2", 1000, "concurrency uppler limit")
	flag.IntVar(&Rounds, "n", 100, "rounds")
	flag.IntVar(&Cmd, "x", CallDefault, "bitwise rpc calls")
	flag.StringVar(&host, "host", "localhost", "rpc server host")
	flag.IntVar(&verbose, "v", 0, "verbose level")
	flag.StringVar(&zk, "zk", "localhost:2181", "zk server addr")
	flag.BoolVar(&testPool, "testpool", false, "test pool")
	flag.BoolVar(&tcpNoDelay, "tcpnodelay", true, "tcp no delay")
	flag.BoolVar(&logTurnOff, "logoff", false, "only show progress instead of rpc result")
	flag.IntVar(&SampleRate, "s", 100, "log sampling rate")
	flag.Usage = showUsage
	flag.Parse()
}

func main() {
	cf := config.NewDefaultProxy()
	cf.PoolCapacity = c2
	cf.IoTimeout = time.Hour
	cf.TcpNoDelay = tcpNoDelay
	proxy := proxy.New(cf)

	etclib.Dial([]string{zk})
	go proxy.StartMonitorCluster()
	proxy.AwaitClusterTopologyReady()

	// test pool
	if testPool {
		testServantPool(proxy)
		pause("pool tested")
	}

	go report.run()

	wg := new(sync.WaitGroup)
	t1 := time.Now()
	for i := 0; i < Rounds; i++ {
		for k := c1; k <= c2; k++ {
			Concurrency = k

			for j := 0; j < k; j++ {
				wg.Add(1)
				go runSession(proxy, wg, i+1, j)
			}

			wg.Wait()
		}

	}

	elapsed := time.Since(t1)
	log.Printf("Elapsed: %s, calls: {%s, %.1f/s}, sessions: {%s, %.1f/s}, errors: {conn:%d, io:%d call:%d}",
		elapsed,
		gofmt.Comma(report.callOk),
		float64(report.callOk)/elapsed.Seconds(),
		gofmt.Comma(int64(report.sessionN)),
		float64(report.sessionN)/elapsed.Seconds(),
		report.connErrs,
		report.ioErrs,
		report.callErrs)
}
