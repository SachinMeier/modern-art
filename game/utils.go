package game

import (
	"crypto/rand"
	"math/big"
)

func randInt(max int) (int, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(i.Int64()), err
}
