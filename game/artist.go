package game

import (
	"fmt"
)

// In order to track points and also allow for tie breakers,
// an ArtPiece is worth 10 points, and a tie breaker is worth 1 point.
// These functions help with that.

const PointsPerArtPiece = 10

// Point returns the number of points for a given number of art pieces
func Point(artPieces int) int {
	return artPieces * PointsPerArtPiece
}

// TieBreakerPoint returns the number of points for a given tie breaker value
func TieBreakerPoint(i int) int {
	return i
}

// Artist is a type for the artists in the game.
type Artist string

const (
	// Manuel = Yellow (12)
	Manuel Artist = "Manuel Carvalho"
	// Sigrid = Blue (13)
	Sigrid Artist = "Sigrid Thaler"
	// Daniel = Red (14)
	Daniel Artist = "Daniel Melim"
	// Ramon = Green (15)
	Ramon Artist = "Ramon Martins"
	// Rafael = Orange (16)
	Rafael Artist = "Rafael Silvera"
	// ArtistNone (non-existent)
	ArtistNone Artist = ""
)

// AllArtists returns a slice of all artists.
func AllArtists() []Artist {
	return []Artist{Manuel, Sigrid, Daniel, Ramon, Rafael}
}

// ArtistArtCounts returns a map of artists to the number of art pieces they have.
func ArtistArtCounts() map[Artist]int {
	return map[Artist]int{
		Manuel: 12,
		Sigrid: 13,
		Daniel: 14,
		Ramon:  15,
		Rafael: 16,
	}
}

// AddTieBreakers adds tiebreaker points to the map of artists.
// Since Artist values are stored as 10 points per ArtPiece in the round,
// the tiebreaker points can never mess up the order.
func AddTieBreakers(artists map[Artist]int) map[Artist]int {
	// return a deep copy to avoid messing with Phase state
	// and allow Phase.RankedArtists() to be called multiple times
	return map[Artist]int{
		Manuel: artists[Manuel] + TieBreakerPoint(4),
		Sigrid: artists[Sigrid] + TieBreakerPoint(3),
		Daniel: artists[Daniel] + TieBreakerPoint(2),
		Ramon:  artists[Ramon] + TieBreakerPoint(1),
		Rafael: artists[Rafael],
	}
}

// ArtPiece is a piece of art, which hails from an Artist.
type ArtPiece struct {
	// Name is an arbitrary name to id the art. necessary?
	Name string
	// Artist is the artist who created the art
	Artist Artist
	// AuctionType
}

// NewArtPiece returns a new ArtPiece with the given name and artist.
func NewArtPiece(artist Artist, name string) ArtPiece {
	return ArtPiece{
		Name:   name,
		Artist: artist,
	}
}

// NewArtPieceDeck returns a slice of all ArtPieces in the game.
// Consider the implementation here. It could be done more lazily
// by tracking count of remaining pieces instead of instantiating
// all of them at once.
func NewArtPieceDeck() []*ArtPiece {
	deck := []*ArtPiece{}
	for artist, count := range ArtistArtCounts() {
		for i := 0; i < count; i++ {
			deck = append(deck, &ArtPiece{
				Name:   fmt.Sprintf("%s-%d", string(artist), i),
				Artist: artist,
			})
		}
	}
	return deck
}

func pickRandomArtPiece(deck []*ArtPiece) (*ArtPiece, []*ArtPiece) {
	idx, err := randInt(len(deck))
	if err != nil {
		panic(err)
	}
	return deck[idx], append(deck[:idx], deck[idx+1:]...)
}
