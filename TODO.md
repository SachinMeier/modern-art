

type GamePlayer struct {
    Name string
    Hand []ArtPiece
    Collection []ArtPiece
    Money int
    Player IPlayer
}


Info a Player should have access to in order to HoldAuction and Bid:
- All Collections
  - How many pieces of each artist are out
  - Who has what
- PastPhases
  - Potential payouts for each artist
- How much money they have
- For single-action auctions, which place they are in the turn order

Things a Player can track on their own but neednt be shared: 
- How many of which ArtPieces were played in each phase
- How much money each other player has


type PhaseNumber int

const (
    Phase1 PhaseNumber = 1
    Phase2 PhaseNumber = 2
    Phase3 PhaseNumber = 3
    Phase4 PhaseNumber = 4
)

type Phase map[Artist]int

Game
- CurrentPhase (PhaseNumber)
- Phases (map[PhaseNumber]Phase)
- Players (list(Player))
- Turn (int)

type Artist string

const (
    Rafael Artist = "Rafael Silvera"
    Sigrid Artist = "Sigrid Thaler"
    Daniel Artist = "Daniel Melim"
    Manuel Artist = "Manuel Carvalho"
    Ramon Artist = "Ramon Martins"
...
)

ArtPiece
- Name/ID (string)
- Artist (Artist)

type Auction struct {
    ArtPiece ArtPiece
    Price int // only used for $ auctions
}

type IPlayer interface {
    HoldAuction() (Auction, error)
    Bid(Auction) (int, error)
}

type Player struct {
    Name string
    Hand []ArtPiece
    Collection []ArtPiece
    Money int
}