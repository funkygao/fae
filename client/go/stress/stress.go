// stress test of fae
package main

import (
	"flag"
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"github.com/funkygao/golib/fixture"
	"github.com/funkygao/golib/gofmt"
	"labix.org/v2/mgo/bson"
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
	clientN     int32
	test1       bool
)

func init() {
	ctx = rpc.NewContext()
	ctx.Caller = "stress test agent"

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func parseFlag() {
	flag.IntVar(&N, "n", 1000, "loops count for each client")
	flag.IntVar(&C, "c", 2500, "concurrent num")
	flag.StringVar(&host, "h", "localhost", "rpc server host")
	flag.IntVar(&verbose, "v", 0, "verbose")
	flag.BoolVar(&test1, "t1", false, "only test connect/close")
	flag.Parse()
}

func sampling() bool {
	return rand.Intn(100) == 1
}

func runClient(proxy *proxy.Proxy, wg *sync.WaitGroup, seq int) {
	defer wg.Done()

	if verbose > 3 && sampling() {
		log.Printf("%6d started\n", seq)
	}

	atomic.AddInt32(&clientN, 1)
	t1 := time.Now()
	client, err := proxy.Servant(host + ":9001")
	if err != nil {
		log.Printf("seq^%d err^%v\n", seq, err)
		atomic.AddInt32(&FailC, 1)
		return
	}
	defer client.Recycle()

	if verbose > 4 && sampling() {
		log.Printf("%6d connected within %s\n", seq, time.Since(t1))
	}

	atomic.AddInt32(&concurrentN, 1)
	defer func() {
		atomic.AddInt32(&concurrentN, -1)
	}()

	var mcKey string
	var mcValue = rpc.NewTMemcacheData()
	var result []byte
	mgQuery, _ := bson.Marshal(bson.M{"snsid": "100003391571259"})
	mgFields, _ := bson.Marshal(bson.M{})
	for i := 0; i < N; i++ {
		_, err = client.Ping(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		mcKey = fmt.Sprintf("mc_stress:%d", rand.Int())
		mcValue.Data = []byte("value of " + mcKey)
		_, err = client.McSet(ctx, mcKey, mcValue, 3600)
		if err != nil {
			log.Println(err)
			return
		}
		_, err, _ = client.McGet(ctx, mcKey)
		_, err = client.LcSet(ctx, mcKey, mcValue.Data)
		if err != nil {
			log.Println(err)
			return
		}
		_, _, err = client.LcGet(ctx, mcKey)
		if err != nil {
			log.Println(err)
			return
		}
		_, _, err = client.IdNext(ctx, 0)
		if err != nil {
			log.Println(err)
			return
		}
		result, _, err = client.MgFindOne(ctx, "default", "idmap", 0,
			mgQuery,
			mgFields)
		if err != nil {
			log.Println(err)
		}

		if false {
			log.Printf("%s", result)
		}
		client.KvdbSet(ctx, fixture.RandomByteSlice(30),
			fixture.RandomByteSlice(10<<10))

		atomic.AddInt64(&calls, 8)
	}

	if verbose > 3 && sampling() {
		log.Printf("%6d done\n", seq)
	}

	atomic.AddInt32(&clientN, -1)
}

func tryServantPool(proxy *proxy.Proxy) {
	for i := 0; i < C; i++ {
		t1 := time.Now()
		client, err := proxy.Servant(host + ":9001")
		if err != nil {
			log.Printf("seq^%d err^%v\n", i, err)
			atomic.AddInt32(&FailC, 1)
			return
		}
		log.Printf("%8d connected within %s", i, time.Since(t1))
		client.Recycle()
	}

	log.Println("try server connect/close done!!!")
}

func main() {
	parseFlag()

	t1 := time.Now()

	proxy := proxy.New(C/2, time.Minute*60)

	tryServantPool(proxy)
	if test1 {
		return
	}

	time.Sleep(time.Second * 2)

	wg := new(sync.WaitGroup)
	for i := 0; i < C; i++ {
		wg.Add(1)
		go runClient(proxy, wg, i)
	}

	// book keeping
	go func() {
		for {
			currentCalls := atomic.LoadInt64(&calls)
			cn := atomic.LoadInt32(&clientN)
			if lastCalls != 0 {
				log.Printf("client: %5d concurrency: %5d calls:%12s cps: %9s\n",
					cn,
					concurrentN,
					gofmt.Comma(currentCalls),
					gofmt.Comma(currentCalls-lastCalls))
			} else {
				log.Printf("client: %5d concurrency: %5d calls: %12s\n",
					cn,
					concurrentN,
					gofmt.Comma(currentCalls))
			}
			lastCalls = currentCalls
			if verbose > 10 {
				log.Println(proxy.StatsJSON())
			}

			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	log.Printf("N=%d, C=%d, FailC=%d, elapsed=%s\n", N, C, FailC, time.Since(t1))
}
