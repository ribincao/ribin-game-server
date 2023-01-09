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
	return nil
}

func (r *TestRoom) AddPlayer(player types.Player) {
	r.playerMap.Store(player.GetId(), player)
}

func TestRoomMethod(t *testing.T) {
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

func TestRoomManager(t *testing.T) {
	room := &TestRoom{
		Id: "test",
	}
	manager.RoomMng.AddRoom(room)
	r := manager.GetRoom[types.Room]("test")
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

func TestPlayerManager(t *testing.T) {
	manager.AddRoomToPlayerMap("xxx", "ribincao")
	manager.AddRoomToPlayerMap("yyy", "whale")
	fmt.Println(manager.GetRoomIdByPlayerId("ribincao"))
	fmt.Println(manager.GetRoomIdByPlayerId("whale"))
}
