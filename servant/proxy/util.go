package proxy

import (
	"io"
	"strings"
)

var (
	ioErrorsFixed = map[string]bool{
		io.EOF.Error():              true,
		io.ErrClosedPipe.Error():    true,
		io.ErrUnexpectedEOF.Error(): true,
	}

	ioErrorsVariant = []string{
		"broken pipe",
		"reset by peer",
	}
)

func IsIoError(err error) bool {
	errmsg := err.Error()
	if _, present := ioErrorsFixed[errmsg]; present {
		return true
	}

	for _, suffix := range ioErrorsVariant {
		if strings.HasSuffix(errmsg, suffix) {
			return true
		}
	}

	return false
}
