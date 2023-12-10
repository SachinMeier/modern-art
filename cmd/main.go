package main

import (
	"encoding/csv"
	"fmt"
	"github.com/SachinMeier/modern-art.git/game"
	"github.com/SachinMeier/modern-art.git/game/players"
	"log"
	"os"
)

func main() {
	simulateAllPossibleCurrentPhases()
	//simulateSinglePhase(0, 0, 0, 0, 0)
	//runManualGame()
}

func runManualGame() {
	p1 := players.NewIOPlayer("me")
	g := game.NewGame([]game.Player{p1})
	g.Start()
}

// NewPhase creates a new phase with the given artist counts
func newPhase(manuel, sigrid, daniel, ramon, rafael int) *game.Phase {
	return &game.Phase{
		ArtistCounts: map[game.Artist]int{
			game.Manuel: game.Point(manuel),
			game.Sigrid: game.Point(sigrid),
			game.Daniel: game.Point(daniel),
			game.Ramon:  game.Point(ramon),
			game.Rafael: game.Point(rafael),
		},
	}
}

func addAuctions(p game.Player, desiredArtistCounts map[game.Artist]int) {
	for artist, count := range desiredArtistCounts {
		for i := 0; i < count; i++ {
			p.HandleAuctionResult(&game.Auction{
				Auctioneer: p,
				WinningBid: game.NewBid(p, 10),
				ArtPiece:   game.NewArtPiece(artist, fmt.Sprintf("added-manuel-%d", i)),
			})
		}
	}
}

func simulateSinglePhase(manuel, sigrid, daniel, ramon, rafael int) {
	p1 := players.NewAlphaPlayer("alpha-1")
	//p2 := players.NewAlphaPlayer("alpha-2")
	hand := []*game.ArtPiece{
		game.NewArtPiece(game.Manuel, "manuel-1"),
		game.NewArtPiece(game.Manuel, "manuel-2"),
		game.NewArtPiece(game.Sigrid, "sigrid-1"),
		game.NewArtPiece(game.Daniel, "daniel-1"),
		game.NewArtPiece(game.Ramon, "ramon-2"),
		game.NewArtPiece(game.Rafael, "rafael-1"),
	}

	// 1 of each artist type
	addAuctions(p1, map[game.Artist]int{
		game.Manuel: manuel,
		game.Sigrid: sigrid,
		game.Daniel: daniel,
		game.Ramon:  ramon,
		game.Rafael: rafael,
	})

	p1.AddArtPieces(hand)
	for _, artist := range game.AllArtists() {
		bid := p1.ExpectedBid(artist)
		log.Printf("expected bid for %s: %d\n", artist, bid)
	}
}

func simulateAllPossibleCurrentPhases() {
	// Create a new CSV file
	file, err := os.Create("scenarios.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"manuel", "sigrid", "daniel", "ramon", "rafael", "manuel-bid", "sigrid-bid", "daniel-bid", "ramon-bid", "rafael-bid",
	}
	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	possiblePieceCts := []int{0, 1, 2, 3, 4}
	for _, manuel := range possiblePieceCts {
		for _, sigrid := range possiblePieceCts {
			for _, daniel := range possiblePieceCts {
				for _, ramon := range possiblePieceCts {
					for _, rafael := range possiblePieceCts {
						p := players.NewAlphaPlayer("alpha-1")
						addAuctions(p, map[game.Artist]int{
							game.Manuel: manuel,
							game.Sigrid: sigrid,
							game.Daniel: daniel,
							game.Ramon:  ramon,
							game.Rafael: rafael,
						})

						strManuel := fmt.Sprintf("%d", manuel)
						strSigrid := fmt.Sprintf("%d", sigrid)
						strDaniel := fmt.Sprintf("%d", daniel)
						strRamon := fmt.Sprintf("%d", ramon)
						strRafael := fmt.Sprintf("%d", rafael)

						bidManuel := fmt.Sprintf("%d", p.ExpectedBid(game.Manuel))
						bidSigrid := fmt.Sprintf("%d", p.ExpectedBid(game.Sigrid))
						bidDaniel := fmt.Sprintf("%d", p.ExpectedBid(game.Daniel))
						bidRamon := fmt.Sprintf("%d", p.ExpectedBid(game.Ramon))
						bidRafael := fmt.Sprintf("%d", p.ExpectedBid(game.Rafael))

						// write to csv
						row := []string{
							strManuel, strSigrid, strDaniel, strRamon, strRafael,
							bidManuel, bidSigrid, bidDaniel, bidRamon, bidRafael,
						}
						if err := writer.Write(row); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}

	// Flush any remaining data in the buffer
	writer.Flush()

	// Check for errors during writing
	if err := writer.Error(); err != nil {
		panic(err)
	}
}
