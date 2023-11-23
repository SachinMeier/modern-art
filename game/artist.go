package game

type Artist string

const (
	// Yellow (12)
	Manuel Artist = "Manuel Carvalho"
	// Blue (13)
	Sigrid Artist = "Sigrid Thaler"
	// Red (14)
	Daniel Artist = "Daniel Melim"
	// Green (15)
	Ramon Artist = "Ramon Martins"
	// Orange (16)
	Rafael Artist = "Rafael Silvera"
)

func AllArtists() []Artist {
	return []Artist{Rafael, Sigrid, Daniel, Manuel, Ramon}
}

// AddTieBreakers adds tie breaker points to the map of artists.
// Since Artist values are stored as 10 points per ArtPiece in the round,
// the tie breaker points can never mess up the order.
func AddTieBreakers(artists map[Artist]int) map[Artist]int {
	artists[Manuel] += 4
	artists[Sigrid] += 3
	artists[Daniel] += 2
	artists[Ramon] += 1
}
