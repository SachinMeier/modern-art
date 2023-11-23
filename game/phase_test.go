package game_test

import (
	"context"
	"github.com/SachinMeier/modern-art.git/game"
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
		phase := game.Phase(map[game.Artist]int{
			game.Manuel: game.Point(1),
			game.Sigrid: game.Point(2),
			game.Daniel: game.Point(3),
			game.Ramon:  game.Point(4),
			game.Rafael: game.Point(5),
		})
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

	// 2. Test that tie breakers are applied correctly
	{
		phase := game.Phase(map[game.Artist]int{
			game.Manuel: game.Point(1),
			game.Sigrid: game.Point(1),
			game.Daniel: game.Point(1),
			game.Ramon:  game.Point(1),
			game.Rafael: game.Point(1),
		})
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

	// 3. Test that tie breakers are applied correctly
	// but don't overrule the points.
	{
		phase := game.Phase(map[game.Artist]int{
			game.Manuel: game.Point(5),
			game.Sigrid: game.Point(1),
			game.Daniel: game.Point(1),
			game.Ramon:  game.Point(2),
			game.Rafael: game.Point(2),
		})
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
}
