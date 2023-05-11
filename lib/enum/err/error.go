package errs

import "errors"

var (
	TOO_LARGE_PACKAGE = errors.New("too large package recv")
	CONNECT_CLOSED    = errors.New("connect closed")
)
