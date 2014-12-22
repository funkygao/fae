package store

// Store of cache
type Store interface {
	Open()
	Close()
	Get(key string) (val interface{}, present bool)
	Put(key string, value interface{})
	Del(key string)
}
