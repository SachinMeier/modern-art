package players_test

import (
	"context"
	"github.com/SachinMeier/modern-art.git/game"
	"github.com/SachinMeier/modern-art.git/game/players"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAlphaPlayerSuite(t *testing.T) {
	suite.Run(t, new(AlphaPlayerTestSuite))
}

type AlphaPlayerTestSuite struct {
	suite.Suite
	testCtx    context.Context
	cancelFunc context.CancelFunc
}

func (suite *AlphaPlayerTestSuite) SetupSuite() {}

func (suite *AlphaPlayerTestSuite) SetupTest() {
	suite.testCtx, suite.cancelFunc = context.WithCancel(context.Background())
}

func (suite *AlphaPlayerTestSuite) TearDownTest() {
	suite.cancelFunc()
}

func (suite *AlphaPlayerTestSuite) TearDownSuite() {}

func (suite *AlphaPlayerTestSuite) Test_HoldAuction() {
	// 1. Test that player returns no auction if he has no art pieces
	{
		p1 := players.NewAlphaPlayer("alpha-1")
		auction, err := p1.HoldAuction()
		suite.Nil(auction)
		suite.ErrorIs(err, game.ErrNoArtPieceToSell)
	}

	// 1. Test that player sells only artist he has
	{
		p1 := players.NewAlphaPlayer("alpha-1")
		m1 := game.NewArtPiece(game.Manuel, "manuel-1")
		p1.AddArtPieces([]*game.ArtPiece{&m1})

		auction, err := p1.HoldAuction()
		if err != nil {
			suite.FailNow("failed to hold auction", err.Error())
		}
		suite.Equal(&m1, auction.ArtPiece)
	}
}
