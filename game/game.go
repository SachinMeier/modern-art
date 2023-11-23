package game

// Game is the main struct for the game.
// It holds all the state of the game.
type Game struct {
	CurrentPhase PhaseNumber
	PastPhases   []Phase
	Players      PlayerOrder
	// ArtPieces is the Deck to be dealt out
	ArtPieces []ArtPiece
}

// NewGame creates a new Game
func NewGame(players []Player) *Game {
	playerOrder := NewPlayerOrder(players)
	g := &Game{
		CurrentPhase: Phase1,
		PastPhases:   []Phase{},
		Players:      playerOrder,
	}
	g.dealArtPieces()
	return g
}

// DoPhase does a phase of the game
func (g *Game) DoPhase() {
	phase := NewPhase()
	for {
		if isOver := g.doTurn(&phase); isOver {
			break
		}
	}
	g.PastPhases = append(g.PastPhases, phase)
	// payout uses CurrentPhase as the index of the phase in PastPhases
	// so we only increment it after payout is done
	g.payout()
	g.CurrentPhase++
	// dealCards uses CurrentPhase to determine how many cards to deal
	g.dealArtPieces()
}

// doTurn does a turn of a Phase
func (g *Game) doTurn(phase *Phase) bool {
	player := g.Players.Pop()

	// Ask the player whose turn it is to hold an auction
	auction, err := player.Player.HoldAuction()
	if err != nil {
		// TODO: handle error? or just crash game
		panic(err)
	}
	phase.AddArtPiece(auction.ArtPiece.Artist)
	// If the auctioned piece ends the round, don't do the auction
	if phase.IsOver() {
		return true
	}

	// push Player to end to allow them to Bid on their own Auction
	// TODO: in the case of a fixed-price auction, Auctioneer bids first, so this
	// should not happen until after the bids are taken in.
	g.Players.Push(player)
	// Copy the player order to allow each player to bid
	auctionBidders := g.Players.Copy()

	// get a bid from each player
	// map[PlayerName]Bid
	// TODO: consider this implementation. To collect data, might be better
	// TODO: to have a map[PlayerName]Bid
	// TODO: here is where we switch-case on AuctionType, for now just 1-shot
	var winningBid Bid
	for _, bidder := range auctionBidders {
		bid, err := bidder.Player.Bid(auction)
		if err != nil {
			panic(err)
		}
		winningBid = BetterBid(winningBid, bid)
	}
	// notify all players of the result
	for _, bidder := range auctionBidders {
		bidder.Player.HandleAuctionResult(winningBid)
	}
	// round can't end after an auction
	return false
}

// payout pays out the players after a concluded Phase
func (g *Game) payout() {
	payouts := g.PastPhases[g.CurrentPhase].Payouts()

	// sum the value of their ArtPiece collection
	for _, player := range g.Players {
		phaseRevenue := 0
		for _, artPiece := range player.Collection {
			phaseRevenue += payouts[artPiece.Artist]
		}
		// give the player the money
		player.Money += phaseRevenue
		// clear colletion
		player.Collection = []ArtPiece{}
	}
}

// TODO: check these numbers
// map[playerCount]map[phaseNumber]artPiecesPerPhase
var artPiecesPerPhase = map[int]map[PhaseNumber]int{
	3: {
		Phase1: 10,
		Phase2: 4,
		Phase3: 4,
	},
	4: {
		Phase1: 8,
		Phase2: 4,
		Phase3: 4,
	},
	5: {
		Phase1: 8,
		Phase2: 3,
		Phase3: 3,
	},
}

func (g *Game) dealArtPieces() {
	piecesToDeal := artPiecesPerPhase[len(g.Players)][g.CurrentPhase]
	for _, player := range g.Players {
		pieces := make([]ArtPiece, piecesToDeal, piecesToDeal)
		for i := 0; i < piecesToDeal; i++ {
			pieces = append(pieces, g.dealArtPiece())
		}
		player.Player.AddArtPieces(pieces)
	}
}

func (g *Game) dealArtPiece() *ArtPiece {
	// TODO: pick a random card
	return nil
}

// Start begins the game
func (g *Game) Start() {
	for range PhaseRange() {
		g.DoPhase()
	}
}
