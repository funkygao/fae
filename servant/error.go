package servant

import (
	"errors"
)

var (
	ErrNotImplemented   = errors.New("Not implemented")
	ErrUnderMaintenance = errors.New("Under maintenance")
)
