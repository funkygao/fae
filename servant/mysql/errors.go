package mysql

import (
	"errors"
)

var (
	ErrNotOpen             = errors.New("mysql: not open")
	ErrServerNotFound      = errors.New("mysql: server not found")
	ErrCircuitOpen         = errors.New("mysql: circuit open")
	ErrShardLookupNotFound = errors.New("mysql: shardId not found in shard lookup table")
)
