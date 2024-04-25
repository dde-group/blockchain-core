package network

type AckWsServerBase struct {
	Event string
	Code  int
	Msg   string
}

const (
	EventError = "error"
	EventSub   = "subbed"
	EventUnSub = "unsub"
)

const (
	OpSubscribe = "sub"
	OpUnSub     = "unsub"
	OpRequest   = "request"
)

const (
	ErrorTooMuch = -1001
)

const (
	ErrMsgTooMuch = "request too much"
)
