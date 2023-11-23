package game

type Player interface {
	HoldAuction() (Auction, error)
	Bid(Auction) (Bid, error)
}
