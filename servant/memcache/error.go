package memcache

import (
	"fmt"
)

type MemcacheError struct {
	Message string
}

func handleError(err *error) {
	if x := recover(); x != nil {
		*err = x.(MemcacheError)
	}
}

func newMemcacheError(format string, args ...interface{}) MemcacheError {
	return MemcacheError{Message: fmt.Sprintf(format, args...)}
}

func (this MemcacheError) Error() string {
	return this.Message
}
