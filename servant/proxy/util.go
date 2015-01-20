package proxy

import (
	"io"
)

var (
	ioErrors = map[string]bool{
		io.EOF.Error():              true,
		io.ErrClosedPipe.Error():    true,
		io.ErrUnexpectedEOF.Error(): true,
		"broken pipe":               true,
	}
)

func IsIoError(err error) bool {
	if _, present := ioErrors[err.Error()]; present {
		return true
	}

	return false
}
