package players

import (
	"errors"
	"github.com/SachinMeier/modern-art.git/game"
	"math"
)

type AlphaPlayer struct {
	name       string
	hand       map[game.Artist][]*game.ArtPiece
	collection []*game.ArtPiece
	money      int

	otherHands []*game.ArtPiece

	currentPhase *game.Phase
	phases       []*game.Phase
	phasePayouts []map[game.Artist]int
}

// Ensures that AlphaPlayer implements game.Player interface at compile time
//var _ game.Player = &AlphaPlayer{}

/*
Logic explainer



*/

// NewAlphaPlayer creates a new AlphaPlayer
func NewAlphaPlayer(name string) *AlphaPlayer {
	return &AlphaPlayer{name: name}
}

// Name returns the Player's name
func (p *AlphaPlayer) Name() string {
	return p.name
}

func (p *AlphaPlayer) SetName(name string) {
	p.name = name
}

func (p *AlphaPlayer) HoldAuction() (*game.Auction, error) {
	artistToSell := game.ArtistNone
	maxExpectedValue := math.MinInt
	for artist, _ := range p.hand {
		expectedValue := p.expectedValue(artist)
		if expectedValue > maxExpectedValue {
			maxExpectedValue = expectedValue
			artistToSell = artist
		}
	}

	if artistToSell == game.ArtistNone {
		return nil, errors.New("no artist to sell")
	}

	return game.NewAuction(p,
		// TODO: for now, sell first art piece of artist. Change when introducing
		// diff auction types
		p.hand[artistToSell][0],
		game.NewBid(p, 0),
	), nil
}

func (p *AlphaPlayer) expectedValue(artist game.Artist) int {
	competitivenessDelta := 0
	selfDelta := 0
	otherDelta := 0
	expectedBid := p.expectedBid(artist)
	return competitivenessDelta*(selfDelta-otherDelta) + expectedBid
}

func (p *AlphaPlayer) expectedBid(artist game.Artist) int {
	// if this art piece would end the round
	if p.artistWouldEndRound(artist) {
		return 0
	}
	return competitiveness * p.averagePayout(artist)
}

func (p *AlphaPlayer) artistWouldEndRound(artist game.Artist) bool {
	return p.currentPhase.ArtistCounts[artist] >= game.MaxArtPiecePointsPerPhase-game.Point(1)
}

func (p *AlphaPlayer) averagePayout(artist game.Artist) int {
	pastPayoutSum := 0
	for _, phasePayout := range p.phasePayouts {
		pastPayoutSum += phasePayout[artist]
	}
	// possible payouts are always 0, pastPayoutSum+10, pastPayoutSum+20, pastPayoutSum+30
	// simplified to 3 * pastPayoutSum + 60
	return (3 * pastPayoutSum) + 60

}

func (p *AlphaPlayer) Bid(auction *game.Auction) (*game.Bid, error) {

}

func (p *AlphaPlayer) HandleAuctionResult(auction *game.Auction) {
	p.currentPhase.AddAuction(auction)
	if p.currentPhase.IsOver() {
		p.phases = append(p.phases, p.currentPhase)
		p.phasePayouts = append(p.phasePayouts, game.CumulativePayouts(p.phases))
		p.currentPhase = game.NewPhase()
	}
}

// AddArtPieces adds ArtPiece's to the Player's Hand
func (p *AlphaPlayer) AddArtPieces(pieces []*game.ArtPiece) {
	for _, piece := range pieces {
		p.hand[piece.Artist] = append(p.hand[piece.Artist], piece)
	}
}

// MoveMoney gives the Player money. Currently only used for payouts.
func (p *AlphaPlayer) MoveMoney(amount int) {
	p.money += amount
}
