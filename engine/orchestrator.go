package engine

import (
	"github.com/coreos/go-etcd/etcd"
	log "github.com/funkygao/log4go"
)

func (this *Engine) startOrchestrator() {
	client := etcd.NewClient(this.conf.etcd.Servers)
	log.Trace("Etcd connected: %+v", *this.conf.etcd)

	watchChan := make(chan *etcd.Response)
	go client.Watch("/hah", 0, false, watchChan, nil)
	for event := range watchChan {
		log.Info("etcd: {action:%s, node:%+v}", event.Action, *event.Node)
	}

	/*
		resp, err := client.Get("frontends", false, false)
		if err != nil {
			log.Fatal(err)
		}

		for _, n := range resp.Node.Nodes {
			log.Printf("%s => %s\n", n.Key, n.Value)
		}

		client.Set("creds", "foo, bar", 0)

		resp, err = client.Get("creds", false, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Current creds: %s: %s\n", resp.Node.Key, resp.Node.Value)
		watchChan := make(chan *etcd.Response)
		go client.Watch("/creds", 0, false, watchChan, nil)
		log.Println("Waiting for an update...")
		r := <-watchChan
		log.Printf("Got updated creds: %s: %s\n", r.Node.Key, r.Node.Value)*/

}
