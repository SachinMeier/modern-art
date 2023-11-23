package game

// GamePlayer is a Player in the Game. It is used to allow the Game
// to keep track of the Player's Hand, Collection, and Money
type GamePlayer struct {
	Player     Player
	Hand       []ArtPiece
	Collection []ArtPiece
	Money      int
}

// NewGamePlayer creates a new GamePlayer from a Player
func NewGamePlayer(player Player) GamePlayer {
	return GamePlayer{
		Player:     player,
		Hand:       []ArtPiece{},
		Collection: []ArtPiece{},
		Money:      0,
	}
}

// TODO: maybe make GamePlayer implement Player interface? use embedding?

// Player is an interface for a player in the Game
type Player interface {
	// Name returns the Player's name
	Name() string
	// HoldAuction requests the Player to put an ArtPiece up for Auction
	HoldAuction() (Auction, error)
	// Bid requests the Player to place a Bid on an Auction
	Bid(Auction) (Bid, error)
	// HandleAuctionResult informs the Player of the result of an Auction by sharing the wining Bid.
	// If a player wins an auction, they should add the ArtPiece to their Collection.
	HandleAuctionResult(Bid)
	// AddArtPieces adds ArtPiece's to the Player's Hand
	AddArtPieces([]ArtPiece)
	// MoveMoney gives the Player money. Currently only used for payouts.
	MoveMoney(int)
}
