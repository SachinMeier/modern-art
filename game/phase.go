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

// PhaseRange returns a slice of all PhaseNumbers
func PhaseRange() []PhaseNumber {
	return []PhaseNumber{Phase1, Phase2, Phase3, Phase4}
}

// Phase contains info about the points for each Artist in a Phase.
type Phase map[Artist]int

// NewPhase creates a new Phase with all artists at 0 points.
func NewPhase() Phase {
	phase := make(map[Artist]int)
	for _, artist := range AllArtists() {
		phase[artist] = 0
	}
	return phase
}

// RankPayouts are the payouts for the top 3 artists in a phase.
const (
	RankPayout1 = 30
	RankPayout2 = 20
	RankPayout3 = 10
)

// MaxArtPiecesPerPhase is the minimum number of art pieces for a given artist
// needed to end the round
const MaxArtPiecesPerPhase = 5

// IsOver returns true if the phase is over.
func (p *Phase) IsOver() bool {
	for _, pieces := range *p {
		// >= allows playing a double when there are 4 pieces down.
		// TODO: check rules on this
		if pieces >= MaxArtPiecesPerPhase {
			return true
		}
	}
	return false
}

// AddArtPiece adds PointsPerArtPiece points to the artist's score.
func (p *Phase) AddArtPiece(artist Artist) {
	(*p)[artist] += PointsPerArtPiece
}

// Winners returns the top 3 artists in the phase.
func (p *Phase) Winners() (Artist, Artist, Artist) {
	artists := p.RankedArtists()
	return artists[0], artists[1], artists[2]
}

// RankedArtists returns a slice of artists sorted by points.
func (p *Phase) RankedArtists() []Artist {
	*p = AddTieBreakers(*p)

	artists := AllArtists()

	sort.SliceStable(artists, func(a, b int) bool {
		return (*p)[artists[a]] > (*p)[artists[b]]
	})

	return artists
}

// Payouts returns a map of artists to their payouts for this Phase
func (p *Phase) Payouts() map[Artist]int {
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
