package servant

import (
	"errors"
)

var (
	ErrNotImplemented    = errors.New("Not implemented")
	ErrInvalidContext    = errors.New("Invalid context")
	ErrServantNotStarted = errors.New("Servant not started")
	ErrMyMergeInvalidRow = errors.New("Row not exist")
)

var (
	ErrProxyNotFound = errors.New("Proxy not found")
)
