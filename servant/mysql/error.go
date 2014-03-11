package mysql

import (
	"errors"
)

var (
	ErrInvalidDsn = errors.New("Invalid mysql dsn")
)
