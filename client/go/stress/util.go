package main

import (
	"log"
	"math/rand"
)

func sampling(rate int) bool {
	return rand.Intn(rate) == 1
}

func showCmdHelp() {
	log.Printf("%16s %3d", "CallPing", CallPing)
	log.Printf("%16s %3d", "CallIdGen", CallIdGen)
	log.Printf("%16s %3d", "CallLCache", CallLCache)
	log.Printf("%16s %3d", "CallMemcache", CallMemcache)
	log.Printf("%16s %3d", "CallMongo", CallMongo)
	log.Printf("%16s %3d", "CallKvdb", CallKvdb)
	log.Printf("%16s %3d", "Ping+Idgen", CallPingIdgen)
	log.Printf("%16s %3d", "Lcache+Idgen", CallIdgenLcache)
}
