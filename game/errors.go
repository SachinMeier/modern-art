package game

import "fmt"

var (
	ErrNotEnoughMoney   = fmt.Errorf("player bid more money than they have")
	ErrArtPieceNotFound = fmt.Errorf("player does not have this card")
	ErrNoArtPieceToSell = fmt.Errorf("player has no art pieces to sell")
)
