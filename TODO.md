
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