// stress test of fae
package main

import (
	"flag"
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/fixture"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	N           int
	C           int
	host        string
	FailC       int32
	ctx         *rpc.Context
	concurrentN int32
	verbose     int
	calls       int64
	lastCalls   int64
)

func init() {
	ctx = rpc.NewContext()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	ctx.Caller = "stress test agent"
}

func parseFlag() {
	flag.IntVar(&N, "n", 10000, "loops count for each client")
	flag.IntVar(&C, "c", 2500, "concurrent num")
	flag.StringVar(&host, "h", "localhost", "rpc server host")
	flag.IntVar(&verbose, "v", 0, "verbose")
	flag.Parse()
}

func runClient(proxy *proxy.Proxy, wg *sync.WaitGroup, seq int) {
	defer wg.Done()

	if verbose > 2 {
		log.Printf("%6d started\n", seq)
	}

	t1 := time.Now()
	client, err := proxy.Servant(host + ":9001")
	if err != nil {
		log.Printf("seq^%d err^%v\n", seq, err)
		atomic.AddInt32(&FailC, 1)
		return
	}
	defer client.Recycle()

	if verbose > 2 {
		log.Printf("%6d connected within %s\n", seq, time.Since(t1))
	}

	atomic.AddInt32(&concurrentN, 1)

	var mcKey string
	var mcValue = rpc.NewTMemcacheData()
	for i := 0; i < N; i++ {
		client.Ping(ctx)
		mcKey = fmt.Sprintf("mc_stress:%d", rand.Int())
		mcValue.Data = []byte("value of " + mcKey)
		client.Ping(ctx)
		client.McAdd(ctx, mcKey, mcValue, 3600)
		client.McSet(ctx, mcKey, mcValue, 3600)
		client.LcSet(ctx, mcKey, mcValue.Data)
		client.LcGet(ctx, mcKey)
		client.IdNext(ctx, 0)
		client.KvdbSet(ctx, fixture.RandomByteSlice(30),
			fixture.RandomByteSlice(10<<10))

		atomic.AddInt64(&calls, 8)
	}

	if verbose > 2 {
		log.Printf("%6d done\n", seq)
	}

	atomic.AddInt32(&concurrentN, -1)
}

func main() {
	parseFlag()

	t1 := time.Now()

	proxy := proxy.New(C, time.Minute*60)
	wg := new(sync.WaitGroup)
	for i := 0; i < C; i++ {
		wg.Add(1)
		go runClient(proxy, wg, i)
	}

	go func() {
		for {
			currentCalls := atomic.LoadInt64(&calls)
			if lastCalls != 0 {
				log.Printf("concurrency: %5d calls:%9d cps:%6d\n", concurrentN,
					atomic.LoadInt64(&calls), currentCalls-lastCalls)
			} else {
				log.Printf("concurrency: %5d calls:%9d\n", concurrentN,
					atomic.LoadInt64(&calls))
			}
			lastCalls = currentCalls
			if verbose > 1 {
				log.Println(proxy.StatsJSON())
			}

			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	log.Printf("N=%d, C=%d, FailC=%d, elapsed=%s\n", N, C, FailC, time.Since(t1))
}
