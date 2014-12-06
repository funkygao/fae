package servant

import (
	"errors"
)

var (
	ErrNotImplemented    = errors.New("Not implemented")
	ErrInvalidContext    = errors.New("Invalid context")
	ErrServantNotStarted = errors.New("Servant not started")
	ErrMyMergeInvalidRow = errors.New("Invalid row to merge")
)
