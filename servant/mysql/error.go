package mysql

import (
	"errors"
)

var (
	ErrServerNotFound = errors.New("mysql: server not found")
)
