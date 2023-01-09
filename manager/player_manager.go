package manager

import "sync"

var player2roomMap sync.Map

func AddRoomToPlayerMap(roomId string, playerId string) {
	if val, ok := player2roomMap.Load(playerId); ok {
		if id, ok := val.(string); ok {
			if id != roomId {
				player2roomMap.Store(playerId, roomId)
			}
		}
	} else {
		player2roomMap.Store(playerId, roomId)
	}
}

func RemoveRoomFromPlayerMap(playerId string) {
	if _, ok := player2roomMap.Load(playerId); ok {
		player2roomMap.Delete(playerId)
	}
}

func GetRoomIdByPlayerId(playerId string) string {
	if val, ok := player2roomMap.Load(playerId); ok {
		if roomId, ok := val.(string); ok {
			return roomId
		}
	}
	return ""
}
