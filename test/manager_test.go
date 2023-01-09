package test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ribincao/ribin-game-server/manager"
	"github.com/ribincao/ribin-game-server/types"
)

type TestRoom struct {
	Id        string
	playerMap sync.Map
}

func (r *TestRoom) GetId() string {
	return r.Id
}

func (r *TestRoom) GetPlayer(playerId string) types.Player {
	player, ok := r.playerMap.Load(playerId)
	if !ok {
		return nil
	}
	return player.(types.Player)
}

func (r *TestRoom) GetAllPlayers() []types.Player {
	var players []types.Player
	r.playerMap.Range(func(key interface{}, value interface{}) bool {
		players = append(players, value.(types.Player))
		return true
	})
	return players
}

func (r *TestRoom) AddPlayer(player types.Player) {
	r.playerMap.Store(player.GetId(), player)
}

func (r *TestRoom) RemovePlayer(playerId string) {
	r.playerMap.Delete(playerId)
}

func TestRoomAddPlayer(t *testing.T) {
	room := &TestRoom{
		Id: "test",
	}
	player := &TestPlayer{
		Id:   "xxx",
		Name: "ribincao",
	}
	room.AddPlayer(player)
	p := room.GetPlayer("xxx")
	if p == nil {
		fmt.Println("player not in room")
		return
	}
	fmt.Println("PlayerName:", p.GetName())
}

func TestRoomDeletePlayer(t *testing.T) {
	room := &TestRoom{
		Id: "test",
	}
	player1 := &TestPlayer{
		Id:   "xxx",
		Name: "ribincao",
	}
	player2 := &TestPlayer{
		Id:   "yyy",
		Name: "whale",
	}
	room.AddPlayer(player1)
	room.AddPlayer(player2)
	p := room.GetPlayer("xxx")
	if p == nil {
		fmt.Println("player not in room")
		return
	}
	ps := room.GetAllPlayers()
	for playerId, player := range ps {
		fmt.Println("before playerId:", playerId, "PlayerName:", player.GetName())
	}

	room.RemovePlayer("yyy")
	ps = room.GetAllPlayers()
	for playerId, player := range ps {
		fmt.Println("after playerId:", playerId, "PlayerName:", player.GetName())
	}
}
func TestRoomManager(t *testing.T) {
	room := &TestRoom{
		Id: "test",
	}
	manager.RoomMng.AddRoom(room)
	r := manager.GetRoom[*TestRoom]("test")
	fmt.Println("playerId:", r.GetId())
}

type TestPlayer struct {
	Id   string
	Name string
}

func (p *TestPlayer) GetId() string {
	return p.Id
}

func (p *TestPlayer) GetName() string {
	return p.Name
}

func (p *TestPlayer) EnterRoom(roomId string) {

}

func TestPlayerManager(t *testing.T) {
	manager.AddRoomToPlayerMap("xxx", "ribincao")
	manager.AddRoomToPlayerMap("yyy", "whale")
	fmt.Println(manager.GetRoomIdByPlayerId("ribincao"))
	fmt.Println(manager.GetRoomIdByPlayerId("whale"))
}
