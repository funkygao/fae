package main

const (
	CallNoop = 1 << iota
	CallPing
	CallIdGen
	CallLCache
	CallMemcache
	CallMongo
	CallMysql
	CallRedis

	CallPingIdgen   = CallPing | CallIdGen
	CallIdgenLcache = CallIdGen | CallLCache
	CallDefault     = CallPing | CallIdGen | CallLCache | CallMysql

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
	zk              string
	testPool        bool
	logTurnOff      bool
	errOff          bool
	tcpNoDelay      bool
	neatStat        bool
	c1              int
	c2              int
)
