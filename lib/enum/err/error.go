package errs

import "errors"

var (
	TOO_LARGE_PACKAGE = errors.New("too large package recv")
	CONNECT_CLOSED    = errors.New("connect closed")
	WORKER_POOL_FULL  = errors.New("worker pool full")
	CONN_NOT_FOUND    = errors.New("connection not found")
	MAX_CONN_REACHED  = errors.New("max connection reached")
)

type Pattern string

const (
	MESSAGE_REGISTERED     Pattern = "msgId: %d has been registered"
	MESSAGE_NOT_REGISTERED Pattern = "msgId: %d has not been registered"
	ACCEPT_TCP_FAILED      Pattern = "Accept TCP failed: %s"
)
