package types

type Player interface {
	GetId() string
	GetName() string
	EnterRoom(roomId string)
}
