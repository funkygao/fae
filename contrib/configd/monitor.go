package main

import (
	"github.com/funkygao/etclib"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

func monitorForever(cf *conf.Conf) {
	loadConfig(cf)
	loadTemplates()

	if err := etclib.Dial(config.etcServers); err != nil {
		panic(err)
	}

	go watchFaes()
	go watchMaintain()

	<-make(chan interface{})
}

func watchFaes() {
	ch := make(chan []string, 10)
	go etclib.WatchService(etclib.SERVICE_FAE, ch)

	for {
		select {
		case <-ch:
			endpoints, err := etclib.ServiceEndpoints(etclib.SERVICE_FAE)
			if err == nil {
				log.Trace("fae endpoints updated: %+v", endpoints)

				dumpFaeConfigPhp(endpoints)
			} else {
				log.Error("fae: %s", err)
			}
		}
	}

	log.Warn("fae watcher died")
}

func watchMaintain() {
	const PATH = "/maintain"

	ch := make(chan []string, 10)
	go etclib.WatchChildren(PATH, ch)

	for {
		select {
		case <-ch:
			kingdoms, err := etclib.Children(PATH)
			if err == nil {
				log.Trace("maintain kingdoms updated: %+v", kingdoms)

				dumpMaintainConfigPhp(kingdoms)
			} else {
				log.Error("maintain kingdom: %s", err)
			}
		}
	}

	log.Warn("maintain watcher died")
}
