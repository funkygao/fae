package redis

import (
	"errors"
)

var (
	ErrCircuitOpen  = errors.New("redis: circuit open")
	ErrPoolNotFound = errors.New("redis pool not found")
	ErrKeyNotExist  = errors.New("key not exists")
)
