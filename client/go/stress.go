// stress test of fae
package main

import (
	"flag"
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"math/rand"
	"sync"
	"time"
)

var (
	N     int
	C     int
	host  string
	FailC int
	ctx   *rpc.Context
)

func init() {
	ctx = rpc.NewContext()
	ctx.Caller = "stress test agent"
}

func parseFlag() {
	flag.IntVar(&N, "n", 10000, "loops count")
	flag.IntVar(&C, "c", 800, "concurrent num")
	flag.StringVar(&host, "h", "localhost", "rpc server host")
	flag.Parse()
}

func runClient(wg *sync.WaitGroup, seq int) {
	defer wg.Done()

	remote := proxy.New()
	client, err := remote.Servant(host + ":9001")
	if err != nil {
		fmt.Printf("seq^%d err^%v\n", seq, err)
		FailC += 1
		return
	}
	defer client.Transport.Close()

	var mcKey string
	var mcValue = rpc.NewTMemcacheData()
	for i := 0; i < N; i++ {
		client.Ping(ctx)
		mcKey = fmt.Sprintf("mc_stress:%d", rand.Int())
		mcValue.Data = []byte("value of " + mcKey)
		client.McAdd(ctx, mcKey, mcValue, 3600)
		client.McSet(ctx, mcKey, mcValue, 3600)
	}
}

func main() {
	parseFlag()

	t1 := time.Now()

	wg := new(sync.WaitGroup)
	for i := 0; i < C; i++ {
		wg.Add(1)
		go runClient(wg, i)
	}
	wg.Wait()

	fmt.Printf("N=%d, C=%d, FailC=%d, elapsed=%s\n", N, C, FailC, time.Since(t1))
}
