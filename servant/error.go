package servant

import (
	"errors"
)

var (
	errLcMissed = errors.New("local cache missed for the key")
)
