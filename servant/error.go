package servant

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
	ErrInvalidContext = errors.New("Invalid context")
)
