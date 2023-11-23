package game

// Auction is an auction for an ArtPiece
type Auction struct {
	// TODO: add incremental id
	// ArtPiece is the piece being auctioned
	ArtPiece *ArtPiece
	// Price is only used in fixed-price Auctions
	// Price int
}

// Bid is a bid on an Auction
type Bid struct {
	PlayerName string
	// TODO: auction id?
	ArtPiece *ArtPiece
	// Value is the amount bid
	Value int
}

// BetterBid returns the better of two Bids, with reigning bid winning ties
func BetterBid(reigning, challenger Bid) Bid {
	if challenger.Value > reigning.Value {
		return challenger
	}
	return reigning
}
