package game

import "fmt"

var (
	ErrNotEnoughMoney   = fmt.Errorf("player bid more money than they have")
	ErrArtPieceNotFound = fmt.Errorf("player does not have this card")
)
