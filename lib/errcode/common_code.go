package errcode

var (
	TOO_LARGE_PACKAGE = NewError(10001, "too large package recv")
	CONNECT_CLOSED    = NewError(10002, "connect closed")
	WORKER_POOL_FULL  = NewError(10003, "worker pool full")
	MAX_CONN_REACHED  = NewError(10004, "max connection reached")
)

var (
	CONN_NOT_FOUND         = NewFormatError(10005, "connection ID %d not found")
	PROPERTY_NOT_FOUND     = NewFormatError(10006, "property %s not found")
	MESSAGE_REGISTERED     = NewFormatError(10007, "msgId: %d has been registered")
	MESSAGE_NOT_REGISTERED = NewFormatError(10008, "msgId: %d has not been registered")
	ACCEPT_TCP_FAILED      = NewFormatError(10009, "Accept TCP failed: %s")
)
