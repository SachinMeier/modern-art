package players

import (
	"crypto/rand"
	"github.com/SachinMeier/modern-art.git/game"
	"math/big"
)

// DummyPlayer is a dummy player for testing
type DummyPlayer struct {
	name       string
	Hand       []*game.ArtPiece
	Collection []*game.ArtPiece
	Money      int
}

// Ensures that DummyPlayer implements Player interface at compile time
var _ game.Player = &DummyPlayer{}

// NewDummyPlayer creates a new DummyPlayer
func NewDummyPlayer(name string) *DummyPlayer {
	return &DummyPlayer{name: name}
}

// Name returns the Player's name
func (dp *DummyPlayer) Name() string {
	return dp.name
}

// SetName sets the Player's name. only used in testing
func (dp *DummyPlayer) SetName(name string) {
	dp.name = name
}

// HoldAuction returns the first card in their Hand
func (dp *DummyPlayer) HoldAuction() (*game.Auction, error) {
	// take first card in Hand
	artPiece := dp.Hand[0]
	// remove it from the Hand
	dp.Hand = dp.Hand[1:]
	return &game.Auction{
		Auctioneer: dp,
		ArtPiece:   artPiece,
	}, nil
}

// Bid requests the Player to place a Bid on an Auction
func (dp *DummyPlayer) Bid(auction *game.Auction) (*game.Bid, error) {
	// bid a random amount up to half of their money
	amount, _ := randInt(dp.Money / 2)
	return &game.Bid{
		Bidder: dp,
		Value:  amount,
	}, nil
}

// HandleAuctionResult informs the Player of the result of an game.Auction
func (dp *DummyPlayer) HandleAuctionResult(auction *game.Auction) {
	if auction.WinningBid.Bidder.Name() == dp.name {
		// add the ArtPiece to their Collection
		dp.Collection = append(dp.Collection, auction.ArtPiece)
	}
}

// AddArtPieces adds ArtPiece's to the Player's Hand
func (dp *DummyPlayer) AddArtPieces(pieces []*game.ArtPiece) {
	dp.Hand = append(dp.Hand, pieces...)
}

// MoveMoney gives the Player money. Currently only used for payouts.
func (dp *DummyPlayer) MoveMoney(amount int) {
	dp.Money += amount
}

func randInt(max int) (int, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(i.Int64()), err
}
