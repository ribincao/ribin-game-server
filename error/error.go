package errs

import "fmt"

const (
	OK                  = 0
	ConfigErrorCode     = 100
	MsgErrorCode        = 200
	ConnectionErrorCode = 300
	TimerErrorCode      = 400

	RoomUnexistErrorCode     = 10000
	PlayerNotInRoomErrorCode = 10001
	ParamsErrorCode          = 10002
)

const (
	FrameworkError = iota
	BusinuessError
)

var (
	ConfigError           = NewFrameworkError(ConfigErrorCode, "config error")
	MsgError              = NewFrameworkError(MsgErrorCode, "msg error")
	ConnectionClosedError = NewFrameworkError(ConnectionErrorCode, "connection closed")
	TimerTickError        = NewFrameworkError(TimerErrorCode, "invalid tick")
	TimerBucketError      = NewFrameworkError(TimerErrorCode, "invalid bucketNum")
	TimerTaskRepeatError  = NewFrameworkError(TimerErrorCode, "task repeat")
	TimerTaskAddError     = NewFrameworkError(TimerErrorCode, "task add error")

	RoomUnexistError     = New(RoomUnexistErrorCode, "room unexist")
	PlayerNotInRoomError = New(PlayerNotInRoomErrorCode, "player not in room")
	RoomIdParamError     = New(ParamsErrorCode, "roomid empty")
	PlayerIdParamError   = New(ParamsErrorCode, "playerid empty")
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
