package game

// GamePlayer is a Player in the Game. It is used to allow the Game
// to keep track of the Player's Hand, Collection, and Money
type GamePlayer struct {
	Player     Player
	Hand       []*ArtPiece
	Collection []*ArtPiece
	Money      int
}

// NewGamePlayer creates a new GamePlayer from a Player
func NewGamePlayer(player Player) *GamePlayer {
	return &GamePlayer{
		Player:     player,
		Hand:       []*ArtPiece{},
		Collection: []*ArtPiece{},
		Money:      0,
	}
}

// RemoveArtPieceFromHand removes an ArtPiece from the GamePlayer's Hand
func (gp *GamePlayer) RemoveArtPieceFromHand(artPiece *ArtPiece) error {
	err := ErrArtPieceNotFound
	for i, piece := range gp.Hand {
		if piece == artPiece {
			gp.Hand = append(gp.Hand[:i], gp.Hand[i+1:]...)
			err = nil
			break
		}
	}
	return err
}

// TODO: maybe make GamePlayer implement Player interface? use embedding?

// Player is an interface for a player in the Game
type Player interface {
	// Name returns the Player's name
	Name() string
	// HoldAuction requests the Player to put an ArtPiece up for Auction
	HoldAuction() (*Auction, error)
	// Bid requests the Player to place a Bid on an Auction
	Bid(*Auction) (*Bid, error)
	// OpenBid requests the Player to place a Bid on an Auction of type AuctionTypeOpen
	// Players should listen for new winning bids on the recv channel and send their bids on the send channel.
	// The game will close the recv channel when the auction is over.
	OpenBid(*Auction, <-chan *Bid, chan<- *Bid)
	// HandleAuctionResult informs the Player of the result of an Auction by sharing the wining Auction.
	// If a player wins an auction, they should add the ArtPiece to their collection.
	HandleAuctionResult(*Auction)
	// AddArtPieces adds ArtPiece's to the Player's hand
	AddArtPieces([]*ArtPiece)
	// MoveMoney gives the Player money. Currently only used for payouts.
	MoveMoney(int)
}
