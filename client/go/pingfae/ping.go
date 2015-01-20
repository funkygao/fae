package main

import (
	"flag"
	"fmt"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	log "github.com/funkygao/log4go"
	"time"
)

const (
	REASON = "pingfae"
)

var (
	host  string
	port  string
	zk    string
	loops int
)

func init() {
	flag.StringVar(&host, "host", "", "host name of faed")
	flag.StringVar(&port, "port", "", "fae port")
	flag.IntVar(&loops, "loops", 1, "ping how many times")
	flag.StringVar(&zk, "zk", "localhost:2181", "zookeeper server addr")
	flag.Parse()

	log.AddFilter("stdout", log.INFO, log.NewConsoleLogWriter())
}

func main() {
	proxy := proxy.NewWithDefaultConfig()
	if host == "" {
		pingCluster(proxy)
		return
	}

	// ping a single faed
	client, err := proxy.ServantByAddr(host + ":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Recycle()

	ctx := rpc.NewContext()
	ctx.Reason = REASON
	ctx.Rid = fmt.Sprintf("req:%d", time.Now().UnixNano())
	pong, err := client.Ping(ctx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pong)
	}

}

func pingCluster(proxy *proxy.Proxy) {
	if err := etclib.Dial([]string{zk}); err != nil {
		fmt.Printf("zk: %s", err.Error())
		return
	}
	go proxy.StartMonitorCluster()
	proxy.AwaitClusterTopologyReady()

	peers := proxy.ClusterPeers()
	if len(peers) == 0 {
		fmt.Println("found no fae peers")
		return
	}

	t := time.Now().Unix() // part of rid
	for i := 0; i < loops; i++ {
		for _, peerAddr := range peers {
			client, err := proxy.ServantByAddr(peerAddr)
			if err != nil {
				fmt.Printf("[%6d] %21s: %s\n", i+1, peerAddr, err)
				continue
			}

			ctx := rpc.NewContext()
			ctx.Reason = REASON
			ctx.Rid = fmt.Sprintf("t:%d,req:%d", t, i+1)
			pong, err := client.Ping(ctx)
			if err != nil {
				client.Close()
				fmt.Printf("[%6d] %21s: %s\n", i+1, peerAddr, err)
			} else {
				fmt.Printf("[%6d] %21s: %s\n", i+1, peerAddr, pong)
			}

			client.Recycle()
		}
	}

}
