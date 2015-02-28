package servant

import (
	"errors"
)

var (
	ErrNotImplemented    = errors.New("Svt: not implemented")
	ErrInvalidContext    = errors.New("Svt: invalid context")
	ErrServantNotStarted = errors.New("Svt: not started")
	ErrMyMergeInvalidRow = errors.New("Svt: row not found")
	ErrProxyNotFound     = errors.New("Svt: proxy not found")
)
