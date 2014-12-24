package store

// Store of cache
type Store interface {
	Get(key string) (val interface{}, present bool)
	Put(key string, val interface{})
	Del(key string)
}
