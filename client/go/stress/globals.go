package main

const (
	CallNoop = 1 << iota
	CallPing
	CallIdGen
	CallLCache
	CallMemcache
	CallMongo
	CallMysql
	CallGame
	CallRedis

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
	logTurnOff      bool
	tcpNoDelay      bool
	c1              int
	c2              int
)
