package memcache

type Hash interface {
	FindServer(key string) *MemcacheClient
	AddServer(addr string)
}
