package store

// Store of cache
type Store interface {
	Get(key string) (val interface{}, present bool)
	Set(key string, val interface{})
	Del(key string)
}
