package game

import "sync"

// Auction is an auction for an ArtPiece
type Auction struct {
	// Auctioneer is the player who is auctioning the ArtPiece
	Auctioneer Player
	// Type is the type of Auction
	Type AuctionType
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

	switch a.Type {
	case AuctionTypeOneShot:
		runOneShotAuction(a, bidders)
	case AuctionTypeOpen:
		runOpenAuction(a, bidders)
	case AuctionTypeBlind:
		runBlindAuction(a, bidders)
	case AuctionTypeSetPrice:
		runSetPriceAuction(a, bidders)
	}
}

// HandleBid compares a new Bid to the current WinningBid and updates the WinningBid if necessary.
// Ties go to the reigning bid.
func (a *Auction) HandleBid(bid *Bid) bool {
	if a.WinningBid == nil {
		a.WinningBid = bid
		return true
	} else if bid.Value > a.WinningBid.Value {
		a.WinningBid = bid
		return true
	}
	return false
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

// TODO: write tests for the different auction types

// AuctionType is the type of Auction
type AuctionType string

// AuctionTypes
const (
	// AuctionTypeOneShot means every player gets one bid sequentially, with the auctioneer going last.
	AuctionTypeOneShot AuctionType = "one-shot"
	// AuctionTypeOpen means every player can submit any number of bids until the auctioneer closes the auction.
	AuctionTypeOpen AuctionType = "open"
	// AuctionTypeBlind means every player submits a single bid simultaneously.
	AuctionTypeBlind AuctionType = "blind"
	// AuctionTypeSetPrice means the auctioneer sets a price and players sequentially get to accept or reject the price.
	AuctionTypeSetPrice AuctionType = "set-price"
)

// runOneShotAuction runs an auction where every player gets one bid sequentially, with the auctioneer going last.
func runOneShotAuction(auction *Auction, bidders []*GamePlayer) {
	// go around and collect bids. Update the auction each time its sent to the next player so they know the current bid
	for _, bidder := range bidders {
		bid, err := bidder.Player.Bid(auction)
		if err != nil {
			panic(err)
		}
		if err := validateBid(bidder, bid); err != nil {
			panic(err)
		}
		auction.HandleBid(bid)
	}
}

// runOpenAuction runs an auction where every player can submit any number of bids until the bidding is done.
// TODO: allow auctioneer to end early? why?
func runOpenAuction(auction *Auction, bidders []*GamePlayer) {
	sends := make([]chan *Bid, len(bidders))
	recvs := make([]chan *Bid, len(bidders))
	// open channel with each player.
	for _, bidder := range bidders {
		// game sends other bids to player
		send := make(chan *Bid, 256)
		// game recv's player's bids
		recv := make(chan *Bid, 256)
		bidder.Player.OpenBid(auction, send, recv)
		sends = append(sends, send)
		recvs = append(recvs, recv)
	}

	bids := make(chan *Bid)
	go funnelBids(bids, recvs)

	for {
		bid, more := <-bids
		if !more {
			// bids channel closed, auction is over. Signal this to all players by closing their recv channels
			for _, recv := range recvs {
				close(recv)
			}
			return
		}
		// if the bid is the new best, tell everyone about it.
		if auction.HandleBid(bid) {
			for _, send := range sends {
				send <- bid
			}
		}
	}
}

// funnelBids funnels bids from many channels into a single channel, for easy synchronous handling
func funnelBids(bids chan *Bid, recvs []chan *Bid) {
	var wg sync.WaitGroup
	wg.Add(len(recvs))
	for _, recv := range recvs {
		go func(recv chan *Bid) {
			for {
				bid, more := <-recv
				// if bidder closed the channel, stop listening
				if !more {
					wg.Done()
					return
				}
				// otherwise, send the bid to the main channel
				bids <- bid
			}
		}(recv)
	}
	// when all bidders have closed their channels, close the bids channel, triggering the end of the auction
	wg.Wait()
	close(bids)
}

// runBlindAuction runs an auction where every player submits a single bid simultaneously.
func runBlindAuction(auction *Auction, bidders []*GamePlayer) {
	// Players do not see one another's bids, so we make a copy of the auction
	// and send each player a copy of the auction with zero starting bid
	staticAuction := NewAuction(auction.Auctioneer, auction.ArtPiece, NewBid(auction.Auctioneer, 0))
	for _, bidder := range bidders {
		bid, err := bidder.Player.Bid(staticAuction)
		if err != nil {
			panic(err)
		}
		if err := validateBid(bidder, bid); err != nil {
			panic(err)
		}
		// the actual auction is updated with the bid
		auction.HandleBid(bid)
	}
}

// runSetPriceAuction runs an auction where the auctioneer sets a price and players sequentially get to accept or reject the price.
func runSetPriceAuction(auction *Auction, bidders []*GamePlayer) {
	// auction starts with auctioneer's bid being the set price. This way, if no one bids, the auctioneer gets the piece
	for _, bidder := range bidders {
		bid, err := bidder.Player.Bid(auction)
		if err != nil {
			panic(err)
		}
		if err := validateBid(bidder, bid); err != nil {
			panic(err)
		}
		if bid.Value == auction.WinningBid.Value {
			// auction is over, price has been accepted
			auction.WinningBid = bid
			return
		}
	}
}
