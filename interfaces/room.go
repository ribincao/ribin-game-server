package interfaces

type Room interface {
	GetId() string
	GetPlayer(playerID string) Player
	GetAllPlayers() []Player
}
