package main

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

const (
	CallPing = 1 << iota
	CallIdGen
	CallLCache
	CallMemcache
	CallMongo
	CallKvdb
	CallMysql
	CallGame

	CallPingIdgen   = CallPing | CallIdGen
	CallIdgenLcache = CallIdGen | CallLCache
	CallDefault     = CallPing | CallIdGen | CallLCache | CallMysql | CallGame

	MC_POOL = "default"
)

var (
	report stats
	ctx    *rpc.Context

	SampleRate      int
	Concurrency     int
	Rounds          int // sessions=Rounds*Concurrency
	LoopsPerSession int // calls=sessions*LoopsPerSession
	Cmd             int
	host            string
	verbose         int
	zk              string
)
