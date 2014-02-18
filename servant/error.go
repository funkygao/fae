package servant

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
	ErrHttp404        = errors.New("Not found")
)
