# Modern Art

Modern Art is a boardgame created by Reiner Knizia. This is a digital version of the game. The original intent was to 
create a flexible implementation allowing for testing different bots, rather than hosting a game for humans to play. That 
is why there is no persistence layer and no web interface. These could likely be added on with minimal effort. 

## Note

The only difference between the boardgame and this implementation I currently know of is the lack of `2x` cards. Any other
difference is likely a bug/oversight.

## Existing Players

The game is meant to be modular and allow different types of players to play together. The game defines a Player interface. 
Any struct that implements the Player interface can be used as a player in the game. Existing Players are defined below.

See the [players/README.md](players/README.md) for more information on existing Player implementations.

## Designing your own Player

You can Implement your own kind of Player quite easily. The interface is defined in `player.go`. 
See the [players/README.md](players/README.md) for more information on how to implement your own player.

### Ideas for Players

1. A player that can be controlled via a web interface. 
2. A player that can be controlled via HTTP API, allowing the actual implementation logic to be run on a different machine or in a different language.
3. A player controlled by an actual AI.
4. A hybrid player that uses a combination of AI and human input to make decisions. The AI might suggest bids/auctions and the Human can confirm or override.

