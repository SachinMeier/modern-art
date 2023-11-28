package main

import (
	"github.com/SachinMeier/modern-art.git/game"
	"github.com/SachinMeier/modern-art.git/game/players"
	"log"
)

func main() {
	//phases := []*game.Phase{
	//	newPhase(1, 1, 3, 4, 5),
	//}

	p1 := players.NewAlphaPlayer("alpha-1")
	p2 := players.NewAlphaPlayer("alpha-2")
	hand := []*game.ArtPiece{
		game.NewArtPiece(game.Manuel, "manuel-1"),
		game.NewArtPiece(game.Manuel, "manuel-2"),
		game.NewArtPiece(game.Sigrid, "sigrid-1"),
		game.NewArtPiece(game.Daniel, "daniel-1"),
		game.NewArtPiece(game.Ramon, "ramon-2"),
		game.NewArtPiece(game.Rafael, "rafael-1"),
	}

	// 1 of each artist type
	{
		p1.HandleAuctionResult(&game.Auction{
			Auctioneer: p2,
			WinningBid: game.NewBid(p2, 10),
			ArtPiece:   game.NewArtPiece(game.Manuel, "other-manuel-1"),
		})
		p1.HandleAuctionResult(&game.Auction{
			Auctioneer: p2,
			WinningBid: game.NewBid(p2, 10),
			ArtPiece:   game.NewArtPiece(game.Manuel, "other-manuel-2"),
		})
		//p1.HandleAuctionResult(&game.Auction{
		//	Auctioneer: p2,
		//	WinningBid: game.NewBid(p2, 10),
		//	ArtPiece:   game.NewArtPiece(game.Manuel, "other-manuel-3"),
		//})
		//p1.HandleAuctionResult(&game.Auction{
		//	Auctioneer: p2,
		//	WinningBid: game.NewBid(p2, 10),
		//	ArtPiece:   game.NewArtPiece(game.Manuel, "other-manuel-4"),
		//})
		p1.HandleAuctionResult(&game.Auction{
			Auctioneer: p2,
			WinningBid: game.NewBid(p2, 10),
			ArtPiece:   game.NewArtPiece(game.Sigrid, "other-sigrid-1"),
		})
		p1.HandleAuctionResult(&game.Auction{
			Auctioneer: p2,
			WinningBid: game.NewBid(p2, 10),
			ArtPiece:   game.NewArtPiece(game.Daniel, "other-daniel-1"),
		})
		p1.HandleAuctionResult(&game.Auction{
			Auctioneer: p2,
			WinningBid: game.NewBid(p2, 10),
			ArtPiece:   game.NewArtPiece(game.Ramon, "other-ramon-1"),
		})
		//p1.HandleAuctionResult(&game.Auction{
		//	Auctioneer: p2,
		//	WinningBid: game.NewBid(p2, 10),
		//	ArtPiece:   game.NewArtPiece(game.Rafael, "other-rafael-1"),
		//})
	}

	p1.AddArtPieces(hand)
	for _, artist := range game.AllArtists() {
		bid := p1.ExpectedBid(artist)
		log.Printf("expected bid for %s: %d\n", artist, bid)
	}

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
