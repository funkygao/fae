package main

import (
	"fmt"
	"github.com/funkygao/fae/engine"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"time"
)

func main() {
	t1 := time.Now()

	client, transport := engine.Servant(":9001")
	defer transport.Close()

	ctx := rpc.NewReqCtx()
	ctx.Caller = "me"
	for i := 0; i < 10; i++ {
		r, _ := client.Ping(ctx)

		fmt.Println(r, time.Since(t1))
		t1 = time.Now()
	}
}
