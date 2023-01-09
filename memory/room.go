package memory

import (
	"sync"

	"github.com/ribincao/ribin-game-server/interfaces"
)

var RoomMng RoomManager

type RoomManager struct {
	RoomMap sync.Map
}

func GetRoom[T interfaces.Room](roomId string) (t T) {
	room := RoomMng.getRoom(roomId)
	if room == nil {
		return
	}
	return room.(T)
}

func (m *RoomManager) getRoom(roomId string) interfaces.Room {
	room, ok := m.RoomMap.Load(roomId)
	if !ok {
		return nil
	}
	return room.(interfaces.Room)
}

func (m *RoomManager) GetRoomIds() []string {
	var roomIds []string
	m.RoomMap.Range(func(key interface{}, value interface{}) bool {
		roomIds = append(roomIds, key.(string))
		return true
	})
	return roomIds
}

func (m *RoomManager) AddRoom(room interfaces.Room) {
	RoomMng.RoomMap.Store(room.GetId(), room)
}

func (m *RoomManager) RemoveRoom(room interfaces.Room) {
	RoomMng.RoomMap.Delete(room.GetId())
}
