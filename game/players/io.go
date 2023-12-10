package players

import (
	"fmt"
	"github.com/SachinMeier/modern-art.git/game"
	"strconv"
)

/*
IOPlayer is a player that requests input from the user via stdin
for each decision. It prints out your hand and collection after
each auction. and shows a log of past events. It is useful for
checking the sanity of other players and for having fun.
*/

// TODO: allow Player to print collection, Other Collections, Money, and Phase Payouts
// with C, O, M, and P respectively

// IOPlayer is a manually controlled player.
type IOPlayer struct {
	name       string
	hand       []*game.ArtPiece
	collection []*game.ArtPiece
	Money      int

	// TODO: possibly make this a map[game.Player][]*game.ArtPiece
	otherCollections map[string][]*game.ArtPiece

	currentPhase *game.Phase
	phases       []*game.Phase
	phasePayouts []map[game.Artist]int
}

// Ensures that IOPlayer implements Player interface at compile time
var _ game.Player = &IOPlayer{}

// NewIOPlayer creates a new IOPlayer
func NewIOPlayer(name string) *IOPlayer {
	return &IOPlayer{
		name:       name,
		hand:       make([]*game.ArtPiece, 0),
		collection: make([]*game.ArtPiece, 0),
		Money:      0,

		otherCollections: make(map[string][]*game.ArtPiece),
		currentPhase:     game.NewPhase(),
		phases:           make([]*game.Phase, 0),
		phasePayouts:     make([]map[game.Artist]int, 0),
	}
}

// Name returns the Player's name
func (p *IOPlayer) Name() string {
	return p.name
}

// SetName sets the Player's name. only used in testing
func (p *IOPlayer) SetName(name string) {
	p.name = name
}

// HoldAuction asks the Player to pick an ArtPiece to auction
func (p *IOPlayer) HoldAuction() (*game.Auction, error) {
	fmt.Printf("Your turn to auction.\n")
	p.printHand()
	fmt.Printf("enter a number to auction that card\n")

	choice := p.handleInput()
	for {
		if choice < 0 || choice >= len(p.hand) {
			fmt.Printf("invalid choice: %d\n", choice)
			// re-request input
			choice = p.handleInput()
			continue
		}
		break
	}

	artPiece := p.hand[choice]
	// remove it from the hand
	p.hand = append(p.hand[:choice], p.hand[choice+1:]...)
	return game.NewAuction(p, artPiece, game.NewBid(p, 0)), nil
}

// Bid requests the Player to place a Bid on an Auction
func (p *IOPlayer) Bid(auction *game.Auction) (*game.Bid, error) {
	fmt.Printf("The following card is up for auction:\n")
	printAuction(auction)
	fmt.Printf("You have %d money. Enter your bid:\n", p.Money)

	bid := p.handleInput()
	for {
		if bid < 0 || bid > p.Money {
			fmt.Printf("invalid bid: %d\n", bid)
			// restart
			bid = p.handleInput()
			continue
		}
		break
	}

	return &game.Bid{
		Bidder: p,
		Value:  bid,
	}, nil
}

// OpenBid requests the Player to place Bid's on an Auction of type AuctionTypeOpen
func (p *IOPlayer) OpenBid(auction *game.Auction, recv <-chan *game.Bid, send chan<- *game.Bid) {
	fmt.Printf("The following card is up for live auction:\n")
	printAuction(auction)
	fmt.Printf("You have %d money. You can keep entering new bids until the auction is over or bid 0 to quit the auction.\n", p.Money)
	// print all new bids
	auctionOver := make(chan bool, 1)
	go func() {
		for {
			bid, more := <-recv
			if !more {
				auctionOver <- true
				return
			}
			fmt.Printf("New best bid: (%s) %d\n", bid.Bidder.Name(), bid.Value)
		}
	}()

	// listen for input
	userBids := make(chan int, 1)

	go func() {
		for {
			choice := p.handleInput()
			if choice == 0 {
				fmt.Printf("You quit the auction. Waiting for auction to end.\n")
				close(send)
				return
			}
			userBids <- choice
		}
	}()

	for {
		select {
		case <-auctionOver:
			fmt.Printf("Auction is over.\n")
			return
		case choice := <-userBids:
			send <- &game.Bid{
				Bidder: p,
				Value:  choice,
			}
		}
	}

}

// HandleAuctionResult informs the Player of the result of a game.Auction
func (p *IOPlayer) HandleAuctionResult(auction *game.Auction) {
	p.currentPhase.AddAuction(auction)
	if p.currentPhase.IsOver() {
		p.phases = append(p.phases, p.currentPhase)
		p.phasePayouts = append(p.phasePayouts, game.CumulativePayouts(p.phases))
		p.currentPhase = game.NewPhase()
	}

	auctionWinner := auction.WinningBid.Bidder.Name()

	if auctionWinner == p.name {
		fmt.Printf("You won the auction!\n")
		fmt.Printf("You paid %d for %s\n", auction.WinningBid.Value, strArtPiece(auction.ArtPiece))
		// add the ArtPiece to their collection
		p.collection = append(p.collection, auction.ArtPiece)
	} else {
		if _, ok := p.otherCollections[auctionWinner]; !ok {
			p.otherCollections[auctionWinner] = make([]*game.ArtPiece, 0)
		}
		p.otherCollections[auctionWinner] = append(p.otherCollections[auctionWinner], auction.ArtPiece)
		fmt.Printf("%s won the auction for %s\n", auctionWinner, strArtPiece(auction.ArtPiece))
	}
}

// AddArtPieces adds ArtPiece's to the Player's hand
func (p *IOPlayer) AddArtPieces(pieces []*game.ArtPiece) {
	p.hand = append(p.hand, pieces...)
	fmt.Printf("You've been dealt new cards:\n")
	p.printHand()
}

// MoveMoney gives the Player money. Currently only used for payouts.
func (p *IOPlayer) MoveMoney(amount int) {
	p.Money += amount
	p.printMoney()
}

func (p *IOPlayer) handleInput() int {
	fmt.Printf("Enter your choice (or enter C, H, O, M, or P):\n")
	var choice string
	if _, err := fmt.Scanln(&choice); err != nil {
		fmt.Printf("error reading input: %s\n", err)
		// restart
		p.handleInput()
	}
	// TODO: use callbacks so that player doesn't forget what they were doing
	switch choice {
	case "C":
		p.printCollection()
		return p.handleInput()
	case "H":
		p.printHand()
		return p.handleInput()
	case "O":
		p.printOtherCollections()
		return p.handleInput()
	case "M":
		p.printMoney()
		return p.handleInput()
	case "P":
		p.printPhasePayouts()
		return p.handleInput()
	// TODO: allow quitting game with Q
	default:
		i, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Printf("invalid choice (must be integer).\n")
			return p.handleInput()
		}
		return i
	}
}

// printers

func (p *IOPlayer) printHand() {
	fmt.Printf("Your hand:\n")
	for i, artPiece := range p.hand {
		fmt.Printf("  %d: %s\n", i, strArtPiece(artPiece))
	}
	printSeparator()
}

func (p *IOPlayer) printCollection() {
	fmt.Printf("Your collection:\n")
	for _, artPiece := range p.collection {
		fmt.Printf("  %s (%s)\n", artPiece.Artist, artPiece.Name)
	}
	printSeparator()
}

func (p *IOPlayer) printOtherCollections() {
	fmt.Printf("Other collections:\n")
	for other, artPieces := range p.otherCollections {
		fmt.Printf("Player %s\n", other)
		for _, artPiece := range artPieces {
			fmt.Printf("  %s (%s)\n", artPiece.Artist, artPiece.Name)
		}
		printSeparator()
	}
}

func (p *IOPlayer) printPhasePayouts() {
	fmt.Printf("Phase payouts:\n")
	for _, artist := range game.AllArtists() {
		fmt.Printf("  %s:", artist)
		for _, phasePayout := range p.phasePayouts {
			fmt.Printf("  %d  ", phasePayout[artist])
		}
		fmt.Println()
	}
	printSeparator()
}

func printSeparator() {
	fmt.Printf("---\n")
}

func strArtPiece(artPiece *game.ArtPiece) string {
	return fmt.Sprintf("%s (%s)", artPiece.Artist, artPiece.Name)
}

func printAuction(auction *game.Auction) {
	fmt.Printf("%s Auction:\n", auction.Type)
	fmt.Printf("  ArtPiece: %s\n", strArtPiece(auction.ArtPiece))
	fmt.Printf("  CurrentBid: %d\n", auction.WinningBid.Value)
	printSeparator()
}

func (p *IOPlayer) printMoney() {
	fmt.Printf("You have %d money\n", p.Money)
}
