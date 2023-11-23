package game

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
	artists[Manuel] += TieBreakerPoint(4)
	artists[Sigrid] += TieBreakerPoint(3)
	artists[Daniel] += TieBreakerPoint(2)
	artists[Ramon] += TieBreakerPoint(1)
	// Rafael doesn't get a tie breaker point

	return artists
}

// ArtPiece is a piece of art, which hails from an Artist.
type ArtPiece struct {
	// Name is an arbitrary name to id the art. necessary?
	Name string
	// Artist is the artist who created the art
	Artist Artist
	// AuctionType
}
