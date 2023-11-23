package game_test

import (
	"context"
	"github.com/SachinMeier/modern-art.git/game"
	"github.com/SachinMeier/modern-art.git/game/players"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPlayerOrderSuite(t *testing.T) {
	suite.Run(t, new(PlayerOrderTestSuite))
}

type PlayerOrderTestSuite struct {
	suite.Suite
	testCtx    context.Context
	cancelFunc context.CancelFunc
}

func (suite *PlayerOrderTestSuite) SetupSuite() {}

func (suite *PlayerOrderTestSuite) SetupTest() {
	suite.testCtx, suite.cancelFunc = context.WithCancel(context.Background())
}

func (suite *PlayerOrderTestSuite) TearDownTest() {
	suite.cancelFunc()
}

func (suite *PlayerOrderTestSuite) TearDownSuite() {}

func (suite *PlayerOrderTestSuite) Test_PlayerOrder_PushPop() {
	// 1. Test that the first player is removed from the PlayerOrder
	{
		dummies := []game.Player{
			players.NewDummyPlayer("1"),
			players.NewDummyPlayer("2"),
			players.NewDummyPlayer("3"),
			players.NewDummyPlayer("4"),
		}

		order := game.NewPlayerOrder(dummies)
		suite.Equal(4, len(order))
		head := order.Pop()
		suite.Equal(3, len(order))
		suite.Equal("1", head.Player.Name())
		order.Push(head)
		suite.Equal(4, len(order))
		head = order.Pop()
		suite.Equal(3, len(order))
		suite.Equal("2", head.Player.Name())
		head = order.Pop()
		suite.Equal(2, len(order))
		suite.Equal("3", head.Player.Name())
		head = order.Pop()
		suite.Equal(1, len(order))
		suite.Equal("4", head.Player.Name())
		// ensure push worked
		head = order.Pop()
		suite.Equal(0, len(order))
		suite.Equal("1", head.Player.Name())

	}
}

func (suite *PlayerOrderTestSuite) Test_PlayerOrder_Copy() {
	// ensure that the copy is a shallow copy
	{
		d1 := players.NewDummyPlayer("1")
		dummies := []game.Player{
			d1,
			players.NewDummyPlayer("2"),
			players.NewDummyPlayer("3"),
			players.NewDummyPlayer("4"),
		}

		order := game.NewPlayerOrder(dummies)

		copyOrder := order.Copy()
		suite.Equal(4, len(copyOrder))

		d1.SetName("A")

		head := copyOrder.Pop()
		suite.Equal("A", head.Player.Name())
		// ensure order is unaffected by Push/Pop
		suite.Equal(4, len(order))
		copyOrder.Push(head)
		suite.Equal(4, len(order))
	}
}
