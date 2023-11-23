package game

import "sort"

type PhaseNumber int

const (
	Phase1 PhaseNumber = iota
	Phase2
	Phase3
	Phase4
)

type Phase map[Artist]int

const PointsPerArtPiece = 10

// NewPhase creates a new Phase with all artists at 0 points.
func NewPhase() Phase {
	phase := make(map[Artist]int)
	for _, artist := range AllArtists() {
		phase[artist] = 0
	}
	return phase
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

func (p *Phase) RankedArtists() []Artist {
	AddTieBreakers(*p)

	artists := AllArtists()

	sort.SliceStable(artists, func(a, b int) bool {
		return (*p)[artists[a]] > (*p)[artists[b]]
	})

	return artists
}
