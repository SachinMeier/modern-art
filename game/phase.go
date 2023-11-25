package game

import "sort"

// PhaseNumber is the number of the phase. 1-4
type PhaseNumber int

// Phase numbers
const (
	Phase1 PhaseNumber = iota
	Phase2
	Phase3
	Phase4
)

// FinalPhase is the last phase. This is used to check if the game is over.
const FinalPhase = Phase4

// PhaseRange returns a slice of all PhaseNumbers
func PhaseRange() []PhaseNumber {
	return []PhaseNumber{Phase1, Phase2, Phase3, Phase4}
}

// Phase contains an ordered list of Auctions
type Phase struct {
	Auctions     []*Auction
	ArtistCounts map[Artist]int
}

// NewPhase creates a new Phase with all artists at 0 points.
func NewPhase() Phase {
	counts := make(map[Artist]int)
	for _, artist := range AllArtists() {
		counts[artist] = 0
	}
	return Phase{
		Auctions:     []*Auction{},
		ArtistCounts: counts,
	}
}

// Len returns the number of Auction in the Phase.
func (p *Phase) Len() int {
	return len(p.Auctions)
}

// RankPayouts are the payouts for the top 3 artists in a phase.
const (
	RankPayout1 = 30
	RankPayout2 = 20
	RankPayout3 = 10
)

// maxArtPiecesPerPhase is the minimum number of art pieces for a given artist
// needed to end the round. DO NOT USE
const maxArtPiecesPerPhase = 2

// MaxArtPiecePointsPerPhase is the minimum number of points for a given artist
// Use this when comparing with ArtistCounts
var MaxArtPiecePointsPerPhase = Point(maxArtPiecesPerPhase)

// IsOver returns true if the given ArtPiece ends the phase.
func (p *Phase) IsOver(artist Artist) bool {
	// >= allows playing a double when there are 4 pieces down.
	// TODO: check rules on this
	// TODO: edit when doubles are possible
	return p.ArtistCounts[artist]+PointsPerArtPiece >= MaxArtPiecePointsPerPhase
}

// AddAuction adds PointsPerArtPiece points to the artist's score
// and appends the Auction to the Phase.
func (p *Phase) AddAuction(auction *Auction) {
	p.Auctions = append(p.Auctions, auction)
	p.ArtistCounts[auction.ArtPiece.Artist] += PointsPerArtPiece
}

// Winners returns the top 3 artists in the phase.
func (p *Phase) Winners() (Artist, Artist, Artist) {
	artists := p.RankedArtists()
	return artists[0], artists[1], artists[2]
}

// RankedArtists returns a slice of artists sorted by points.
func (p *Phase) RankedArtists() []Artist {
	artistCounts := AddTieBreakers(p.ArtistCounts)

	artists := AllArtists()

	sort.SliceStable(artists, func(a, b int) bool {
		return artistCounts[artists[a]] > artistCounts[artists[b]]
	})

	return artists
}

// phasePayouts returns a map of artists to their payouts for this isolated Phase
func (p *Phase) phasePayouts() map[Artist]int {
	artists := p.RankedArtists()
	payouts := make(map[Artist]int)
	for i, artist := range artists {
		switch i {
		case 0:
			payouts[artist] = RankPayout1
		case 1:
			payouts[artist] = RankPayout2
		case 2:
			payouts[artist] = RankPayout3
		default:
			payouts[artist] = 0
		}
	}

	return payouts
}

// CumulativePayouts returns a map of artists to their cumulative payouts
// given a list of phases
func CumulativePayouts(phases []Phase) map[Artist]int {
	prevPayouts := make(map[Artist]int)
	lastIdx := len(phases) - 1
	// sum all but the most recent phase
	for _, phase := range phases[:lastIdx] {
		for artist, payout := range phase.phasePayouts() {
			prevPayouts[artist] += payout
		}
	}

	lastPhase := phases[lastIdx].phasePayouts()
	for artist, payout := range lastPhase {
		if lastPhase[artist] != 0 {
			lastPhase[artist] = payout + prevPayouts[artist]
		}
	}

	return lastPhase
}
