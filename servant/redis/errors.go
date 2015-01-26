package redis

import (
	"errors"
)

var (
	ErrDataNotExists = errors.New("key not exists")
	ErrCircuitOpen   = errors.New("redis: circuit open")
	ErrPoolNotFound  = errors.New("redis pool not found")
	ErrKeyNotExist   = errors.New("key not exists")
)
