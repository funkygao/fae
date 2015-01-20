package proxy

import (
	"errors"
	"github.com/funkygao/assert"
	"testing"
)

func TestIsIoError(t *testing.T) {
	err := errors.New("EOF")
	assert.Equal(t, true, IsIoError(err))
	err = errors.New("broken pipe")
	assert.Equal(t, true, IsIoError(err))
}
