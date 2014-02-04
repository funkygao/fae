package mongo

import (
	"errors"
)

var (
	ErrServerNotFound = errors.New("mongodb server not found")
)
