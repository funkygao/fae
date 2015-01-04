package redis

import (
	"errors"
)

var (
	ErrorDataNotExists = errors.New("key not exists")
	ErrCircuitOpen     = errors.New("redis: circuit open")
)
