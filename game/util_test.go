package game_test

import (
	"github.com/SachinMeier/modern-art.git/game"
	"github.com/stretchr/testify/suite"
)

func newArtPiece(artist game.Artist) game.ArtPiece {
	return game.NewArtPiece(artist, "test")
}

func newAuctionWithWinningBid(auctioneer game.Player, artist game.Artist,
	bidder game.Player, value int) *game.Auction {
	artPiece := newArtPiece(artist)
	return game.NewAuction(auctioneer, &artPiece, game.NewBid(bidder, value))
}

// Matchers

func mustMatchAuctions(suite *suite.Suite, a1, a2 *game.Auction) {
	mustMatchPlayers(suite, a1.Auctioneer, a2.Auctioneer)
	mustMatchArtPiece(suite, a1.ArtPiece, a2.ArtPiece)
	mustMatchBids(suite, a1.WinningBid, a2.WinningBid)
}

func mustMatchBids(suite *suite.Suite, bid1, bid2 *game.Bid) {
	mustMatchPlayers(suite, bid1.Bidder, bid2.Bidder)
	suite.Equal(bid1.Value, bid2.Value)
}

func mustMatchPlayers(suite *suite.Suite, p1, p2 game.Player) {
	suite.Equal(p1.Name(), p2.Name())
}

func mustMatchGamePlayers(suite *suite.Suite, p1, p2 *game.GamePlayer) {
	suite.Equal(p1.Player.Name(), p2.Player.Name())
	suite.Equal(p1.Money, p2.Money)
	suite.Equal(p1.Collection, p2.Collection)
	suite.Equal(p1.Hand, p2.Hand)
}

func mustMatchArtPiece(suite *suite.Suite, a1, a2 *game.ArtPiece) {
	suite.Equal(a1.Name, a2.Name)
	suite.Equal(a1.Artist, a2.Artist)
}
