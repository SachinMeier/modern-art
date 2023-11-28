package players

import (
	"github.com/SachinMeier/modern-art.git/game"
)

type ArtistCount struct {
	Artist game.Artist
	Count  int
}

func rankedArtistCounts(phase *game.Phase) []ArtistCount {
	artistCounts := make([]ArtistCount, 0, len(phase.ArtistCounts))

	for _, artist := range phase.RankedArtists() {
		artistCounts = append(artistCounts, ArtistCount{artist, phase.ArtistCounts[artist]})
	}

	return artistCounts
}
