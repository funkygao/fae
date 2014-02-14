// stress test of fae
package main

import (
	"flag"
	"fmt"
	"github.com/funkygao/fae/proxy"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"sync"
	"time"
)

var (
	N   int
	C   int
	ctx *rpc.Context
)

func init() {
	ctx = rpc.NewContext()
	ctx.Caller = "stress test agent"
}

func parseFlag() {
	flag.IntVar(&N, "n", 10000, "loops count")
	flag.IntVar(&C, "c", 500, "concurrent num")
	flag.Parse()
}

func runClient(wg *sync.WaitGroup) {
	defer wg.Done()

	remote := proxy.New()
	client, err := remote.Servant(":9001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Transport.Close()

	for i := 0; i < N; i++ {
		client.Ping(ctx)
	}
}

func main() {
	parseFlag()

	t1 := time.Now()

	wg := new(sync.WaitGroup)
	for i := 0; i < C; i++ {
		wg.Add(1)
		go runClient(wg)
	}
	wg.Wait()

	fmt.Printf("N=%d, elapsed=%s", N, time.Since(t1))
}
