package game_test

import (
	"context"
	"github.com/SachinMeier/modern-art.git/game"
	"github.com/SachinMeier/modern-art.git/game/players"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPhaseSuite(t *testing.T) {
	suite.Run(t, new(PhaseTestSuite))
}

type PhaseTestSuite struct {
	suite.Suite
	testCtx    context.Context
	cancelFunc context.CancelFunc
}

func (suite *PhaseTestSuite) SetupSuite() {}

func (suite *PhaseTestSuite) SetupTest() {
	suite.testCtx, suite.cancelFunc = context.WithCancel(context.Background())
}

func (suite *PhaseTestSuite) TearDownTest() {
	suite.cancelFunc()
}

func (suite *PhaseTestSuite) TearDownSuite() {}

func (suite *PhaseTestSuite) Test_PhaseRankedArtists() {
	// 1. Test that the artists are sorted by points.
	{
		phase := game.Phase{
			ArtistCounts: map[game.Artist]int{
				game.Manuel: game.Point(1),
				game.Sigrid: game.Point(2),
				game.Daniel: game.Point(3),
				game.Ramon:  game.Point(4),
				game.Rafael: game.Point(5),
			},
		}
		order := []game.Artist{
			game.Rafael,
			game.Ramon,
			game.Daniel,
			game.Sigrid,
			game.Manuel,
		}

		artists := phase.RankedArtists()

		for i, artist := range artists {
			suite.Equal(order[i], artist)
		}
	}

	// 2. Test that tiebreakers are applied correctly
	{
		phase := game.Phase{
			ArtistCounts: map[game.Artist]int{
				game.Manuel: game.Point(1),
				game.Sigrid: game.Point(1),
				game.Daniel: game.Point(1),
				game.Ramon:  game.Point(1),
				game.Rafael: game.Point(1),
			},
		}
		order := []game.Artist{
			game.Manuel,
			game.Sigrid,
			game.Daniel,
			game.Ramon,
			game.Rafael,
		}

		artists := phase.RankedArtists()

		for i, artist := range artists {
			suite.Equal(order[i], artist)
		}
	}

	// 3. Test that tiebreakers are applied correctly
	// but don't overrule the points.
	{
		phase := game.Phase{
			ArtistCounts: map[game.Artist]int{
				game.Manuel: game.Point(5),
				game.Sigrid: game.Point(1),
				game.Daniel: game.Point(1),
				game.Ramon:  game.Point(2),
				game.Rafael: game.Point(2),
			},
		}
		order := []game.Artist{
			game.Manuel,
			game.Ramon,
			game.Rafael,
			game.Sigrid,
			game.Daniel,
		}

		artists := phase.RankedArtists()

		for i, artist := range artists {
			suite.Equal(order[i], artist)
		}
	}

	// 4. Test empty phase
	{
		phase := game.Phase{
			ArtistCounts: map[game.Artist]int{},
		}
		order := []game.Artist{
			game.Manuel,
			game.Sigrid,
			game.Daniel,
			game.Ramon,
			game.Rafael,
		}
		artists := phase.RankedArtists()

		for i, artist := range artists {
			suite.Equal(order[i], artist)
		}
	}
}

func (suite *PhaseTestSuite) Test_PhaseAddAuction() {
	// 1. Test that the winning bid is added to the phase
	{
		p1 := players.NewDummyPlayer("1")
		p2 := players.NewDummyPlayer("2")
		phase := game.NewPhase()

		auction1 := newAuctionWithWinningBid(p1, game.Manuel, p2, 1)
		phase.AddAuction(auction1)
		suite.Equal(1, len(phase.Auctions))
		mustMatchAuctions(&suite.Suite, auction1, phase.Auctions[0])
		suite.Equal(game.Point(1), phase.ArtistCounts[game.Manuel])

		auction2 := newAuctionWithWinningBid(p2, game.Sigrid, p1, 2)
		phase.AddAuction(auction2)
		suite.Equal(2, len(phase.Auctions))
		mustMatchAuctions(&suite.Suite, auction2, phase.Auctions[1])
		suite.Equal(game.Point(1), phase.ArtistCounts[game.Sigrid])

		auction3 := newAuctionWithWinningBid(p1, game.Manuel, p2, 3)
		phase.AddAuction(auction3)
		suite.Equal(3, len(phase.Auctions))
		mustMatchAuctions(&suite.Suite, auction3, phase.Auctions[2])
		suite.Equal(game.Point(2), phase.ArtistCounts[game.Manuel])

		suite.Equal(game.Point(2), phase.ArtistCounts[game.Manuel])
		suite.Equal(game.Point(1), phase.ArtistCounts[game.Sigrid])
		suite.Equal(game.Point(0), phase.ArtistCounts[game.Daniel])
		suite.Equal(game.Point(0), phase.ArtistCounts[game.Ramon])
		suite.Equal(game.Point(0), phase.ArtistCounts[game.Rafael])
	}
}

func (suite *PhaseTestSuite) Test_PhaseIsOver() {
	// 1. New Phase is not over
	{
		phase := game.NewPhase()
		suite.False(phase.IsOver())
	}

	// 2. Phase is over based on artist
	{
		phase := newPhase(0, 1, 2, 3, 5)
		suite.True(phase.IsOver())
	}
}

func (suite *PhaseTestSuite) Test_PhasePayouts() {
	// 1. Test that the payouts are correct for single phase
	{
		phase := newPhase(0, 1, 2, 3, 5)
		payouts := phase.PhasePayouts()
		suite.Equal(5, len(payouts))

		suite.Equal(0, payouts[game.Manuel])
		suite.Equal(0, payouts[game.Sigrid])
		suite.Equal(game.RankPayout3, payouts[game.Daniel])
		suite.Equal(game.RankPayout2, payouts[game.Ramon])
		suite.Equal(game.RankPayout1, payouts[game.Rafael])
	}

	// only 2 winners
	{
		phase := newPhase(0, 0, 0, 3, 5)
		payouts := phase.PhasePayouts()
		suite.Equal(5, len(payouts))

		suite.Equal(0, payouts[game.Manuel])
		suite.Equal(0, payouts[game.Sigrid])
		suite.Equal(0, payouts[game.Daniel])
		suite.Equal(game.RankPayout2, payouts[game.Ramon])
		suite.Equal(game.RankPayout1, payouts[game.Rafael])
	}

	// only 1 winner
	{
		phase := newPhase(5, 0, 0, 0, 0)
		payouts := phase.PhasePayouts()
		suite.Equal(5, len(payouts))

		suite.Equal(game.RankPayout1, payouts[game.Manuel])
		suite.Equal(0, payouts[game.Sigrid])
		suite.Equal(0, payouts[game.Daniel])
		suite.Equal(0, payouts[game.Ramon])
		suite.Equal(0, payouts[game.Rafael])
	}

	// test tiebreaker scenario
	{
		phase := newPhase(5, 1, 1, 1, 0)
		payouts := phase.PhasePayouts()
		suite.Equal(5, len(payouts))

		suite.Equal(game.RankPayout1, payouts[game.Manuel])
		suite.Equal(game.RankPayout2, payouts[game.Sigrid])
		suite.Equal(game.RankPayout3, payouts[game.Daniel])
		suite.Equal(0, payouts[game.Ramon])
		suite.Equal(0, payouts[game.Rafael])
	}
}

func (suite *PhaseTestSuite) Test_CumulativePayouts() {
	// test 1 round
	{
		phases := []*game.Phase{
			newPhase(1, 2, 3, 4, 5),
		}
		payouts := game.CumulativePayouts(phases)
		suite.Equal(5, len(payouts))
		suite.Equal(0, payouts[game.Manuel])
		suite.Equal(0, payouts[game.Sigrid])
		suite.Equal(game.RankPayout3, payouts[game.Daniel])
		suite.Equal(game.RankPayout2, payouts[game.Ramon])
		suite.Equal(game.RankPayout1, payouts[game.Rafael])
	}

	// test 2 rounds
	{
		phases := []game.Phase{
			newPhase(1, 2, 3, 4, 5),
			newPhase(5, 4, 3, 2, 1),
		}
		payouts := game.CumulativePayouts(phases)
		suite.Equal(5, len(payouts))
		suite.Equal(game.RankPayout1, payouts[game.Manuel])
		suite.Equal(game.RankPayout2, payouts[game.Sigrid])
		suite.Equal(game.RankPayout3+game.RankPayout3, payouts[game.Daniel])
		suite.Equal(0, payouts[game.Ramon])
		suite.Equal(0, payouts[game.Rafael])
	}

	// test 3 rounds
	{
		phases := []game.Phase{
			newPhase(1, 2, 3, 4, 5),
			newPhase(5, 4, 3, 2, 1),
			newPhase(0, 0, 5, 2, 3),
		}
		payouts := game.CumulativePayouts(phases)
		suite.Equal(5, len(payouts))
		suite.Equal(0, payouts[game.Manuel])
		suite.Equal(0, payouts[game.Sigrid])
		suite.Equal(game.RankPayout1+game.RankPayout3+game.RankPayout3, payouts[game.Daniel])
		suite.Equal(game.RankPayout3+game.RankPayout2, payouts[game.Ramon])
		suite.Equal(game.RankPayout2+game.RankPayout1, payouts[game.Rafael])
	}

	// test 4 rounds, past ones incomplete
	{
		phases := []game.Phase{
			newPhase(0, 0, 0, 5, 4),
			newPhase(5, 0, 0, 0, 0),
			newPhase(0, 0, 5, 0, 0),
			newPhase(1, 2, 3, 4, 5),
		}
		payouts := game.CumulativePayouts(phases)
		suite.Equal(5, len(payouts))
		suite.Equal(0, payouts[game.Manuel])
		suite.Equal(0, payouts[game.Sigrid])
		suite.Equal(game.RankPayout3+game.RankPayout1, payouts[game.Daniel])
		suite.Equal(game.RankPayout2+game.RankPayout1, payouts[game.Ramon])
		suite.Equal(game.RankPayout1+game.RankPayout2, payouts[game.Rafael])
	}

	// test 4 rounds, current one incomplete
	{
		phases := []game.Phase{
			newPhase(0, 0, 0, 5, 4),
			newPhase(5, 0, 0, 0, 0),
			newPhase(0, 0, 5, 0, 0),
			newPhase(0, 2, 5, 0, 0),
		}
		payouts := game.CumulativePayouts(phases)
		suite.Equal(5, len(payouts))
		suite.Equal(0, payouts[game.Manuel])
		suite.Equal(game.RankPayout2, payouts[game.Sigrid])
		suite.Equal(game.RankPayout1+game.RankPayout1, payouts[game.Daniel])
		suite.Equal(0, payouts[game.Ramon])
		suite.Equal(0, payouts[game.Rafael])
	}

}
