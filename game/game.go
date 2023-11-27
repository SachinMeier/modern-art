package game

const (
	MaxPlayers    = 5
	StartingMoney = 100
)

// Game is the main struct for the game.
// It holds all the state of the game.
type Game struct {
	CurrentPhase PhaseNumber
	PastPhases   []*Phase
	Players      PlayerOrder
	// ArtPieces is the Deck to be dealt out
	ArtPieces []*ArtPiece
}

// NewGame creates a new Game
func NewGame(players []Player) *Game {
	if len(players) > MaxPlayers {
		panic("too many players")
	}
	playerOrder := NewPlayerOrder(players)
	g := &Game{
		CurrentPhase: Phase1,
		PastPhases:   []*Phase{},
		Players:      playerOrder,
		ArtPieces:    NewArtPieceDeck(),
	}

	for _, player := range g.Players {
		g.givePlayerMoney(player, StartingMoney)
	}

	return g
}

// Start begins the game
func (g *Game) Start() map[string]int {
	for {
		if gameOver := g.DoPhase(); gameOver {
			break
		}
	}
	return g.CalculateScores()
}

// DoPhase does a phase of the game. Returns true if game is over
func (g *Game) DoPhase() bool {
	// dealCards uses CurrentPhase to determine how many cards to deal
	g.DealArtPieces()
	phase := NewPhase()
	for {
		if isOver := g.doTurn(phase); isOver {
			break
		}
	}
	g.PastPhases = append(g.PastPhases, phase)
	// PayoutPlayer uses CurrentPhase as the index of the phase in PastPhases
	// so we only increment it after payout is done
	g.PayoutPlayers()
	if gameOver := g.NextPhase(); gameOver {
		return true
	}
	return false
}

// doTurn does a turn of a Phase
func (g *Game) doTurn(phase *Phase) bool {
	auctioneer := g.Players.Pop()

	// Ask the auctioneer whose turn it is to hold an auction
	auction, err := auctioneer.Player.HoldAuction()
	if err != nil {
		// TODO: handle error? or just crash game
		panic(err)
	}
	// TODO: validate auction and maybe even set auction.Auctioneer
	// remove the ArtPiece from the auctioneer's hand
	if err := auctioneer.RemoveArtPieceFromHand(auction.ArtPiece); err != nil {
		panic(err)
	}
	// If the auctioned piece ends the round, don't do the auction
	phase.AddAuction(auction)
	if phase.IsOver() {
		// set auction Bid to nil to indicate no winner. This is necessary
		// for fixed-price auctions where the auctioneer bids first. Then
		// add it to the phase to allow for payouts & scoring
		auction.WinningBid = nil
		// we need to push the auctioneer back on the end of the queue
		// to ensure all auctioneers are remembered for payouts & scoring
		g.Players.Push(auctioneer)
		return true
	}

	// push Player to end to allow them to Bid on their own Auction
	// TODO: in the case of a fixed-price auction, Auctioneer bids first, so this
	// should not happen until after the bids are taken in.
	g.Players.Push(auctioneer)

	// Copy the auctioneer order to allow each player to bid
	auctionBidders := g.Players.Copy()
	auction.Run(auctionBidders)

	// notify all auctioneers of the result
	for _, bidder := range auctionBidders {
		bidder.Player.HandleAuctionResult(auction)
	}

	buyer := g.LookupGamePlayer(auction.WinningBid.Bidder.Name())
	// debit money from the buyer
	g.givePlayerMoney(buyer, -auction.WinningBid.Value)
	// give buyer the art piece
	buyer.Collection = append(buyer.Collection, auction.ArtPiece)
	// if the auctioneer bought their own painting, money goes to the bank
	// else give money to the auctioneer
	if auctioneer.Player.Name() != auction.WinningBid.Bidder.Name() {
		g.givePlayerMoney(auctioneer, auction.WinningBid.Value)
	}
	// round can't end after an auction
	return false
}

// NextPhase increments the CurrentPhase
func (g *Game) NextPhase() bool {
	g.CurrentPhase++
	return g.GameOver()
}

// LookupGamePlayer returns the GamePlayer for a given Player name
// NOTE: returns nil, leading to a panic, if the player is not found
// TODO: Possibly add map of players to make this O(1), but its max O(5) anyway
func (g *Game) LookupGamePlayer(name string) *GamePlayer {
	for _, player := range g.Players {
		if player.Player.Name() == name {
			return player
		}
	}
	return nil
}

// GameOver returns true if the game is over
func (g *Game) GameOver() bool {
	return g.CurrentPhase > FinalPhase
}

// PayoutPlayers pays out the players after a concluded Phase
func (g *Game) PayoutPlayers() {
	payouts := CumulativePayouts(g.PastPhases)

	// sum the value of their ArtPiece collection
	for _, player := range g.Players {
		phaseRevenue := 0
		for _, artPiece := range player.Collection {
			phaseRevenue += payouts[artPiece.Artist]
		}
		// give the player the money
		g.givePlayerMoney(player, phaseRevenue)
		// clear collection
		player.Collection = []*ArtPiece{}
	}
}

// ArtPiecesPerPhase map[playerCount]map[phaseNumber]artPiecesPerPhase
// TODO: check these numbers
var ArtPiecesPerPhase = map[int]map[PhaseNumber]int{
	3: {
		Phase1: 10,
		Phase2: 4,
		Phase3: 4,
		Phase4: 0,
	},
	4: {
		Phase1: 8,
		Phase2: 4,
		Phase3: 4,
		Phase4: 0,
	},
	5: {
		Phase1: 8,
		Phase2: 3,
		Phase3: 3,
		Phase4: 0,
	},
}

// DealArtPieces deals ArtPieces to the players
func (g *Game) DealArtPieces() {
	piecesToDeal := ArtPiecesPerPhase[len(g.Players)][g.CurrentPhase]
	for _, player := range g.Players {
		pieces := make([]*ArtPiece, piecesToDeal, piecesToDeal)
		for i := 0; i < piecesToDeal; i++ {
			pieces[i] = g.dealArtPiece()
		}
		g.givePlayerArtPieces(player, pieces)
	}
}

func (g *Game) dealArtPiece() *ArtPiece {
	artPiece, remaining := pickRandomArtPiece(g.ArtPieces)
	g.ArtPieces = remaining
	return artPiece
}

func (g *Game) givePlayerMoney(player *GamePlayer, amount int) {
	player.Money += amount
	// notify player of payout
	player.Player.MoveMoney(amount)
}

func (g *Game) givePlayerArtPieces(player *GamePlayer, artPieces []*ArtPiece) {
	player.Hand = append(player.Hand, artPieces...)
	// notify player of new art pieces
	player.Player.AddArtPieces(artPieces)
}

// CalculateScores returns a map of player names -> scores
func (g *Game) CalculateScores() map[string]int {
	scores := make(map[string]int)
	for _, player := range g.Players {
		scores[player.Player.Name()] = player.Money
	}
	return scores
}
