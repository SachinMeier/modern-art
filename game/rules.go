package game

const (
	MaxPlayers    = 5
	StartingMoney = 100
)

// RankPayouts are the payouts for the top 3 artists in a phase.
const (
	RankPayout1 = 30
	RankPayout2 = 20
	RankPayout3 = 10
)

// MaxArtPiecesPerPhase is the minimum number of art pieces for a given artist
// needed to end the round.
// 5
const MaxArtPiecesPerPhase = 5

// FinalPhase is the last phase. This is used to check if the game is over.
const FinalPhase = Phase4

// Tiebreakers are the tiebreaker points for each artist.
var Tiebreakers = map[Artist]int{
	Manuel: TieBreakerPoint(4),
	Sigrid: TieBreakerPoint(3),
	Daniel: TieBreakerPoint(2),
	Ramon:  TieBreakerPoint(1),
	Rafael: TieBreakerPoint(0),
}

// ArtPiecesPerPhase map[playerCount]map[phaseNumber]artPiecesPerPhase
// TODO: check these numbers
var ArtPiecesPerPhase = map[int]map[PhaseNumber]int{
	1: {
		Phase1: 3,
		Phase2: 2,
		Phase3: 2,
		Phase4: 0,
	},
	2: {
		Phase1: 12,
		Phase2: 8,
		Phase3: 8,
		Phase4: 0,
	},
	3: {
		Phase1: 10,
		Phase2: 6,
		Phase3: 6,
		Phase4: 0,
	},
	4: {
		Phase1: 8,
		Phase2: 4,
		Phase3: 4,
		Phase4: 0,
	},
	5: {
		Phase1: 8,
		Phase2: 3,
		Phase3: 3,
		Phase4: 0,
	},
}
