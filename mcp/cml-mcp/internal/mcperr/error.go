package mcperr

type Error struct {
	Code int
	Msg  string
	Data any
}

func New(message string) *Error {
	return &Error{Code: -32000, Msg: message}
}

func WithCode(message string, code int) *Error {
	return &Error{Code: code, Msg: message}
}

func WithData(message string, data any) *Error {
	return &Error{Code: -32000, Msg: message, Data: data}
}

func (e *Error) Error() string {
	return e.Msg
}
