package game

// Auction is an auction for an ArtPiece
type Auction struct {
	Auctioneer Player
	// TODO: add incremental id
	// ArtPiece is the piece being auctioned
	ArtPiece *ArtPiece
	// WinningBid is the winning Bid for the Auction
	WinningBid *Bid
}

// NewAuction creates a new Auction
func NewAuction(auctioneer Player, artPiece *ArtPiece, winningBid *Bid) *Auction {
	return &Auction{
		Auctioneer: auctioneer,
		ArtPiece:   artPiece,
		WinningBid: winningBid,
	}
}

// Run runs an Auction by requesting a single bid from each bidder.
// TODO: later, this is where we will add support for different AuctionType's
func (a *Auction) Run(bidders []*GamePlayer) {
	// get a bid from each player
	// map[PlayerName]Bid
	// TODO: consider this implementation. To collect data, might be better
	// TODO: to have a map[PlayerName]Bid
	// TODO: here is where we switch-case on AuctionType, for now just 1-shot
	for _, bidder := range bidders {
		bid, err := bidder.Player.Bid(a)
		if err != nil {
			panic(err)
		}
		if err := validateBid(bidder, bid); err != nil {
			panic(err)
		}
		a.HandleBid(bid)
	}
}

// HandleBid compares a new Bid to the current WinningBid and updates the WinningBid if necessary.
// Ties go to the reigning bid.
func (a *Auction) HandleBid(bid *Bid) {
	if a.WinningBid == nil {
		a.WinningBid = bid
	} else if bid.Value > a.WinningBid.Value {
		a.WinningBid = bid
	}
}

func validateBid(player *GamePlayer, bid *Bid) error {
	if bid.Value > player.Money {
		return ErrNotEnoughMoney
	}
	return nil
}

// Bid is a bid on an Auction
type Bid struct {
	Bidder Player
	// Value is the amount bid
	Value int
}

// NewBid creates a new Bid
func NewBid(bidder Player, value int) *Bid {
	return &Bid{
		Bidder: bidder,
		Value:  value,
	}
}
