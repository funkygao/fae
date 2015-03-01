// test a broken fae connection
// HOW to play:
// run broken while faed is on, turn off fae, turn on fae
// see output of broken
// TODO try zk cluster
package main

import (
	"flag"
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	log "github.com/funkygao/log4go"
	"time"
)

var (
	host     string
	port     string
	interval int
)

func init() {
	flag.StringVar(&host, "host", "localhost", "host name of faed")
	flag.StringVar(&port, "port", "9001", "fae port")
	flag.IntVar(&interval, "interval", 5, "sleep between conn")
	flag.Parse()

	log.AddFilter("stdout", log.INFO, log.NewConsoleLogWriter())
}

func main() {
	cf := config.NewDefaultProxy()
	cf.PoolCapacity = 2
	peer := proxy.New(cf)

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		fmt.Println()

		client, err := peer.ServantByAddr(host + ":" + port)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ctx := rpc.NewContext()
		ctx.Reason = "broken.test"
		ctx.Rid = time.Now().UnixNano()
		pong, err := client.Ping(ctx)
		if err != nil {
			fmt.Println("err:", err)

			if proxy.IsIoError(err) {
				client.Close()
			}
		} else {
			fmt.Println(pong)
		}

		client.Recycle()
	}

}
