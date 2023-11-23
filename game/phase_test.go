package game_test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"modern-art/game"
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

func (suite *PhaseTestSuite) TestPhaseRankedArtists() {
	phase := map[game.Artist]int{
		game.Manuel: 1 * game.PointsPerArtPiece,
		game.Sigrid: 2 * game.PointsPerArtPiece,
		game.Daniel: 3 * game.PointsPerArtPiece,
		game.Ramon:  4 * game.PointsPerArtPiece,
		game.Rafael: 5 * game.PointsPerArtPiece,
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
