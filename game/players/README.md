# Players

This directory holds several implementations of the Player interface found in `../game/player.go`. 
Each one is a different strategy for playing the game.

### Dummy Player

File: `dummy.go`

The dummy player has the dumbest logic possible but does properly play the game. It is meant to be a template
and for testing the game itself. I don't suggest playing the game with Dummy players.

### IO Player

File: `io.go`

The IO Player takes input from the command line and outputs to the command line. It allows a human to play the game
as if it were a text-based game.

### Alpha Player

File: `alpha.go`

Alpha Player is my first attempt at creating a fixed-algorithm player. I implemented a rough approximation of the logic
I use when evaluating decisions in the game.

Currently, Alpha considers how many cards have been played, the number of cards played for each artist, and the tiebreakers.
It does not consider the following factors which I think a future version should:
- The number of cards per artist left in the Hands/Deck
- Who plays next and what they are incentivized to play
- How playing a specific Artist would benefit current collections of self and other players