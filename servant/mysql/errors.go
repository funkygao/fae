package mysql

import (
	"errors"
)

var (
	ErrNotOpen        = errors.New("mysql: not open")
	ErrServerNotFound = errors.New("mysql: server not found")
	ErrCircuitOpen    = errors.New("mysql: circuit open")
)
