package players

import (
	"fmt"
	"github.com/SachinMeier/modern-art.git/game"
	"log"
	"math"
)

/*
AlphaPlayer is my first attempt at an AI player. It is a simple
and will most likely always play the artist in first place.
It will bid pretty close to the expected value of the artist.
*/
type AlphaPlayer struct {
	name       string
	hand       map[game.Artist][]*game.ArtPiece
	collection map[game.Artist][]*game.ArtPiece
	money      int

	// TODO: possibly make this a map[game.Player][]*game.ArtPiece
	otherHands [][]*game.ArtPiece

	currentPhase *game.Phase
	phases       []*game.Phase
	phasePayouts []map[game.Artist]int
}

// Ensures that AlphaPlayer implements game.Player interface at compile time
var _ game.Player = &AlphaPlayer{}

// NewAlphaPlayer creates a new AlphaPlayer
func NewAlphaPlayer(name string) *AlphaPlayer {
	return &AlphaPlayer{
		name:       name,
		hand:       make(map[game.Artist][]*game.ArtPiece),
		collection: make(map[game.Artist][]*game.ArtPiece),
		money:      0,

		otherHands:   make([][]*game.ArtPiece, 0),
		currentPhase: game.NewPhase(),
		phases:       make([]*game.Phase, 0),
		phasePayouts: make([]map[game.Artist]int, 0),
	}
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
		expectedValue := p.ExpectedValue(artist)
		if expectedValue > maxExpectedValue {
			maxExpectedValue = expectedValue
			artistToSell = artist
		}
	}

	if artistToSell == game.ArtistNone {
		return nil, game.ErrNoArtPieceToSell
	}

	return game.NewAuction(p,
		// TODO: for now, sell first art piece of artist. Change when introducing
		// diff auction types
		p.hand[artistToSell][0],
		game.NewBid(p, 0),
	), nil
}

func (p *AlphaPlayer) ExpectedValue(artist game.Artist) int {
	competitivenessDelta := 0
	selfDelta := 0
	otherDelta := 0
	expectedBid := p.ExpectedBid(artist)
	return competitivenessDelta*(selfDelta-otherDelta) + expectedBid
}

func (p *AlphaPlayer) ExpectedBid(artist game.Artist) int {
	// if this art piece would end the round
	if p.artistWouldEndRound(artist) {
		return 0
	}
	competitiveness := p.calculateCompetitiveness(artist)
	log.Printf("competitiveness for %s: %f\n", artist, competitiveness)
	// normalize the competitiveness scale from 100 to the payout range
	maxPayout := p.maxPayout(artist)
	scale := float64(maxPayout) / 100.0

	expectedBid := int(math.Floor(competitiveness * scale))
	return nonNegative(expectedBid)
}

func nonNegative(i int) int {
	if i < 0 {
		return 0
	}
	return i
}

const oneThird = 1.0 / 3.0

// calculateCompetitiveness calculates how likely an artist is to place.
// aiming for 100 as guaranteed first place, 0 as guaranteed last place (latter being possible)
func (p *AlphaPlayer) calculateCompetitiveness(artist game.Artist) float64 {
	/*
		simple logic is that base case is each card 1-5 is worth 20 points.
		because its 1/5th of the way towards first place. We add the tiebreakers
		to give a slight edge to tiebreaker winners.

		Then, we look at the competition and add or deduct smaller amounts
		based on how much more or less we have than the other artists.
	*/
	n := p.currentPhase.ArtistCounts[artist]

	artPieceBaseFactor := int(100.0 / float64(game.MaxArtPiecePointsPerPhase))
	competitiveness := float64(n*artPieceBaseFactor + game.Tiebreakers[artist])

	// since we're considering playing this artist, it will get a boost of
	// one more art piece
	//newPieceBaseFactor := float64(game.MaxArtPiecePointsPerPhase) / 10.0
	// scales down how much each comparison matters
	pieceScaleFactor := 2.0
	for i, ct := range rankedArtistCounts(p.currentPhase) {
		// skip self
		if artist == ct.Artist {
			continue
		}

		// how much more does other artist have than self
		// divide by Points to see how many pcs diff between self and other
		nLead := float64(n+game.PointsPerArtPiece-ct.Count) / float64(game.MaxArtPiecePointsPerPhase)

		// simpler algorithm
		//if nLead > game.PointsPerArtPiece {
		//	// self is already ahead of this guy.
		//	competitiveness += placeWeights(i)
		//} else if nLead < 0 && nLead > -game.PointsPerArtPiece {
		//	// if self is played, it will overtake this artist
		//	competitiveness += placeWeights(i)/2.0
		//} else {
		//	// if self is played, it will not overtake this artist
		//	competitiveness += 0
		//}

		// take cube root of nLead to get a diminishing return
		// in either direction.
		cubedDelta := cubeRoot(nLead)
		log.Printf("cubedDelta for %s vs. %s: %f\n", artist, ct.Artist, cubedDelta)
		weightedDelta := placeWeights(i, cubedDelta)

		// placedWeights makes competitiveness more related
		// to the first place than the third place (for example)
		// TODO: some constant is needed here
		competitiveness = competitiveness + (weightedDelta * pieceScaleFactor)
	}
	return competitiveness
}

func cubeRoot(i float64) float64 {
	// since math.Pow doesn't accept negatives, we need to
	// make it positive if negative, take the cube root, then negate it
	negate := 1.0
	if i < 0 {
		negate = -1.0
		i = -1 * i
	}
	cubedDelta := math.Pow(i, oneThird)
	return negate * cubedDelta
}

// TODO: this should be based on whether cubedDelta is positive or negative
// beating a first place artist shouldbe significant, while losing to it should not,
// and vice versa for last place
func placeWeights(i int, val float64) float64 {
	// TODO: make this a const?
	artistCt := len(game.AllArtists())
	// losing to 4th place is as significant as beating 1st place
	if val < 0 {
		i = artistCt - 1 - i
	}
	// TODO: make this a continuous function, in the math sense
	switch i {
	// should never happen: losing to last and beating 1st are impossible
	case 0:
		return 1.2 * val
	// beat second place or lose to 4th place
	case 1:
		return 1.0 * val
	// beat or lose to third place
	case 2:
		return 0.7 * val
	// lose to second place or beat fourth place
	case 3:
		return 0.45 * val
	// lose to first place or beat last place
	case 4:
		return 0.2 * val
	// should never happen
	default:
		return 0.0
	}
}

func (p *AlphaPlayer) artistWouldEndRound(artist game.Artist) bool {
	return p.currentPhase.ArtistCounts[artist] >= game.MaxArtPiecePointsPerPhase-game.Point(1)
}

func (p *AlphaPlayer) pastPayoutSum(artist game.Artist) int {
	pastPayoutSum := 0
	for _, phasePayout := range p.phasePayouts {
		pastPayoutSum += phasePayout[artist]
	}
	return pastPayoutSum
}

func (p *AlphaPlayer) maxPayout(artist game.Artist) int {
	pastPayoutSum := p.pastPayoutSum(artist)
	return pastPayoutSum + game.RankPayout1
}

func (p *AlphaPlayer) averagePayout(artist game.Artist) float64 {
	pastPayoutSum := p.pastPayoutSum(artist)
	// possible payouts are always 0, pastPayoutSum+10, pastPayoutSum+20, pastPayoutSum+30
	// simplified to 3 * pastPayoutSum + 60
	return float64(3*pastPayoutSum+rankPayoutSum()) / 4.0
}

func rankPayoutSum() int {
	return game.RankPayout1 + game.RankPayout2 + game.RankPayout3
}

func (p *AlphaPlayer) Bid(auction *game.Auction) (*game.Bid, error) {
	return nil, fmt.Errorf("not implemented")
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
		if val, ok := p.hand[piece.Artist]; !ok {
			p.hand[piece.Artist] = []*game.ArtPiece{piece}
		} else {
			p.hand[piece.Artist] = append(val, piece)
		}
	}
}

// MoveMoney gives the Player money. Currently only used for payouts.
func (p *AlphaPlayer) MoveMoney(amount int) {
	p.money += amount
}
