package main

const (
	CallPing = 1 << iota
	CallIdGen
	CallLCache
	CallMemcache
	CallMongo
	CallMysql
	CallGame

	CallPingIdgen   = CallPing | CallIdGen
	CallIdgenLcache = CallIdGen | CallLCache
	CallDefault     = CallPing | CallIdGen | CallLCache | CallMysql | CallGame

	MC_POOL = "default"
)

var (
	report stats

	SampleRate      int
	Concurrency     int
	Rounds          int // sessions=Rounds*Concurrency
	LoopsPerSession int // calls=sessions*LoopsPerSession
	Cmd             int
	host            string
	verbose         int
	zk              string
	testPool        bool
)
