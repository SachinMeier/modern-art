package game_test

import (
	"context"
	"fmt"
	"github.com/SachinMeier/modern-art.git/game"
	"github.com/SachinMeier/modern-art.git/game/players"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}

type GameTestSuite struct {
	suite.Suite
	testCtx    context.Context
	cancelFunc context.CancelFunc
}

func (suite *GameTestSuite) SetupSuite() {}

func (suite *GameTestSuite) SetupTest() {
	suite.testCtx, suite.cancelFunc = context.WithCancel(context.Background())
}

func (suite *GameTestSuite) TearDownTest() {
	suite.cancelFunc()
}

func (suite *GameTestSuite) TearDownSuite() {}

func (suite *GameTestSuite) Test_DealArtPieces() {
	// 1. Test that cards are dealt to players correctly each phase
	{
		playerCt := 4
		dummies := suite.getNDummyPlayers(playerCt)
		ng := game.NewGame(dummies)
		ng.DealArtPieces()
		p1ct := game.ArtPiecesPerPhase[playerCt][game.Phase1]
		for _, gp := range ng.Players {
			suite.Equal(p1ct, len(gp.Hand))
		}

		ng.NextPhase()
		ng.DealArtPieces()
		p2ct := game.ArtPiecesPerPhase[playerCt][game.Phase2]
		for _, gp := range ng.Players {
			suite.Equal(p1ct+p2ct, len(gp.Hand))
		}

		ng.NextPhase()
		ng.DealArtPieces()
		p3ct := game.ArtPiecesPerPhase[playerCt][game.Phase3]
		for _, gp := range ng.Players {
			suite.Equal(p1ct+p2ct+p3ct, len(gp.Hand))
		}

		ng.NextPhase()
		ng.DealArtPieces()
		p4ct := game.ArtPiecesPerPhase[playerCt][game.Phase4]
		for _, gp := range ng.Players {
			suite.Equal(p1ct+p2ct+p3ct+p4ct, len(gp.Hand))
		}
	}
}

func (suite *GameTestSuite) Test_GameOver() {
	// 1. Test that the game ends when the last phase is over
	playerCt := 4
	dummies := suite.getNDummyPlayers(playerCt)
	ng := game.NewGame(dummies)
	for range game.PhaseRange() {
		suite.False(ng.GameOver())
		ng.NextPhase()
	}
	suite.True(ng.GameOver())
}

func (suite *GameTestSuite) Test_Game() {
	playerCt := 4
	dummies := suite.getNDummyPlayers(playerCt)
	ng := game.NewGame(dummies)

	// check that the players have the correct amount of money
	for _, player := range ng.Players {
		suite.Equal(game.StartingMoney, player.Money)
	}
	// do phase 1
	isGameOver := ng.DoPhase()
	suite.False(isGameOver)

	// check number of artPieces played
	phase1 := ng.PastPhases[0]
	phase1ArtCt := phase1.Len()
	sumOfArtInHands := 0
	for _, player := range ng.Players {
		sumOfArtInHands += len(player.Hand)
	}
	// artPieces dealt = artPieces played + artPieces in hands
	suite.Equal(playerCt*game.ArtPiecesPerPhase[playerCt][game.Phase1], phase1ArtCt+sumOfArtInHands,
		"artPieces dealt (%d) = played (%d) + in hands (%d)", playerCt*game.ArtPiecesPerPhase[playerCt][game.Phase1],
		phase1ArtCt, sumOfArtInHands)

	// check there is one winner for the round
	winnerCt := 0
	for _, ct := range ng.PastPhases[0].ArtistCounts {
		if ct >= game.MaxArtPiecePointsPerPhase {
			winnerCt += 1
		}
	}
	suite.Equal(1, winnerCt)

	// check that the players all got the correct amount of money
	p1First, p1Second, p1Third := phase1.Winners()
	playerMoney := ng.CalculateScores()
	expectedMoney := map[string]int{}
	for _, dummy := range dummies {
		expectedMoney[dummy.Name()] = game.StartingMoney
	}

	for _, auction := range phase1.Auctions[:phase1ArtCt-1] {
		// should panic if auction.winningBid is nil because that should only be true
		// for the last auction
		auctioneer := auction.Auctioneer.Name()
		buyer := auction.WinningBid.Bidder.Name()
		// add what auctioneer got
		if buyer != auctioneer {
			expectedMoney[auctioneer] += auction.WinningBid.Value
		}
		// subtract what buyer paid
		expectedMoney[buyer] -= auction.WinningBid.Value
		// add what payout gave
		switch auction.ArtPiece.Artist {
		case p1First:
			expectedMoney[buyer] += game.RankPayout1
		case p1Second:
			expectedMoney[buyer] += game.RankPayout2
		case p1Third:
			expectedMoney[buyer] += game.RankPayout3
		}
	}

	for name, money := range expectedMoney {
		suite.Equal(money, playerMoney[name], "incorrect payout for %s", name)
	}
	// do phase 2

}

// helpers

func (suite *GameTestSuite) getNDummyPlayers(n int) []game.Player {
	ps := make([]game.Player, n, n)
	for i := 0; i < n; i++ {
		ps[i] = players.NewDummyPlayer(fmt.Sprintf("dummy-%d", i))
	}
	return ps
}
