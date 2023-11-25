package game

// PlayerOrder a FIFO queue of players
type PlayerOrder []*GamePlayer

// NewPlayerOrder creates a new PlayerOrder from a slice of Players
func NewPlayerOrder(players []Player) PlayerOrder {
	gamePlayers := make([]*GamePlayer, len(players), len(players))
	for i, player := range players {
		gamePlayers[i] = NewGamePlayer(player)
	}
	return gamePlayers
}

// Pop removes the first player from the PlayerOrder and returns it
func (po *PlayerOrder) Pop() *GamePlayer {
	player := (*po)[0]
	*po = (*po)[1:]
	return player
}

// Push adds a player to the end of the PlayerOrder,
// should be called at the end of their turn
func (po *PlayerOrder) Push(player *GamePlayer) {
	*po = append(*po, player)
}

// Copy returns a shallow copy of the PlayerOrder for iterating through an Auction
func (po *PlayerOrder) Copy() PlayerOrder {
	newPo := make(PlayerOrder, len(*po))
	copy(newPo, *po)
	return newPo
}
