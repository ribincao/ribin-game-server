package types

type Room interface {
	GetId() string
	GetPlayer(playerID string) Player
	GetAllPlayers() []Player
	AddPlayer(Player)
	RemovePlayer(playerId string)
}
