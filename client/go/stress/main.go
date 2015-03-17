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
	flag.IntVar(&c1, "c1", 10, "concurrency lower limit")
	flag.IntVar(&c2, "c2", 5000, "concurrency uppler limit")
	flag.IntVar(&Rounds, "n", 100, "rounds")
	flag.IntVar(&Cmd, "x", CallDefault, "bitwise rpc calls")
	flag.StringVar(&host, "host", "localhost", "rpc server host")
	flag.StringVar(&zk, "zk", "localhost:2181", "zk server addr")
	flag.BoolVar(&testPool, "testpool", false, "test pool")
	flag.BoolVar(&tcpNoDelay, "tcpnodelay", true, "tcp no delay")
	flag.BoolVar(&logTurnOff, "logoff", false, "only show progress instead of rpc result")
	flag.BoolVar(&errOff, "erroff", false, "turn off err log")
	flag.BoolVar(&neatStat, "neatstat", true, "only show concurrency and qps in stats output")
	flag.IntVar(&SampleRate, "s", 100, "log sampling rate")
	flag.Usage = showUsage
	flag.Parse()
}

func main() {
	cf := config.NewDefaultProxy()
	cf.IoTimeout = time.Hour
	cf.TcpNoDelay = tcpNoDelay
	prx := proxy.New(cf)

	etclib.Dial([]string{zk})
	go prx.StartMonitorCluster()
	prx.AwaitClusterTopologyReady()

	// test pool
	if testPool {
		testServantPool(prx)
		pause("pool tested")
	}

	go report.run()

	wg := new(sync.WaitGroup)
	t1 := time.Now()
	for k := c1; k <= c2; k += 10 {
		Concurrency = k

		cf.PoolCapacity = Concurrency
		prx = proxy.New(cf)

		for i := 0; i < Rounds; i++ {
			for j := 0; j < k; j++ {
				wg.Add(1)
				go runSession(prx, wg, i+1, j)
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
