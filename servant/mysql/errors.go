package mysql

import (
	"errors"
)

var (
	ErrNotOpen             = errors.New("mysql not open")
	ErrServerNotFound      = errors.New("mysql server not found")
	ErrCircuitOpen         = errors.New("mysql circuit open")
	ErrShardLookupNotFound = errors.New("shard lookup fails")
	ErrLookupTableNotFound = errors.New("mysql lookup table not configured")
)
