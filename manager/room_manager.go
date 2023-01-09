package manager

import (
	"sync"

	"github.com/ribincao/ribin-game-server/types"
)

var RoomMng RoomManager

type RoomManager struct {
	RoomMap sync.Map
}

func GetRoom[T types.Room](roomId string) (t T) {
	room := RoomMng.getRoom(roomId)
	if room == nil {
		return
	}
	return room.(T)
}

func (m *RoomManager) getRoom(roomId string) types.Room {
	room, ok := m.RoomMap.Load(roomId)
	if !ok {
		return nil
	}
	return room.(types.Room)
}

func (m *RoomManager) GetRoomIds() []string {
	var roomIds []string
	m.RoomMap.Range(func(key interface{}, value interface{}) bool {
		roomIds = append(roomIds, key.(string))
		return true
	})
	return roomIds
}

func (m *RoomManager) AddRoom(room types.Room) {
	RoomMng.RoomMap.Store(room.GetId(), room)
}

func (m *RoomManager) RemoveRoom(room types.Room) {
	RoomMng.RoomMap.Delete(room.GetId())
}
