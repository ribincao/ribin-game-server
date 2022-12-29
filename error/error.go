package errs

import "fmt"

const (
	OK                  = 0
	ConfigErrorCode     = 100
	MsgErrorCode        = 200
	ConnectionErrorCode = 300
)

const (
	FrameworkError = iota
	BusinuessError
)

var (
	ConfigError           = NewFrameworkError(ConfigErrorCode, "config error")
	MsgError              = NewFrameworkError(MsgErrorCode, "msg error")
	ConnectionClosedError = NewFrameworkError(ConnectionErrorCode, "connection closed")
)

type Error struct {
	error
	Code    int32
	Type    int
	Message string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Type == FrameworkError {
		return fmt.Sprintf("type : framework, code : %d, msg : %s", e.Code, e.Message)
	}
	return fmt.Sprintf("type : business, code : %d, msg : %s", e.Code, e.Message)
}

func NewFrameworkError(code int32, msg string) *Error {
	return &Error{
		Type:    FrameworkError,
		Code:    code,
		Message: msg,
	}
}

func New(code int32, msg string) *Error {
	return &Error{
		Type:    BusinuessError,
		Code:    code,
		Message: msg,
	}
}
