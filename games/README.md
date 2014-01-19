Repository of games supported by Battle of Bits

## Writing a Game

Add a folder to this project. Inside the folder, add your Go source code.

Your game should have a `Match` object with the following interface:

```go
type Match interface {

	// Retrieve the player whose turn it is.
	CurrentPlayer() *Player

	// Serialize a move from the byte string and apply it to the board. Returns
	// `true` if the game is over, and an error if the move was unreadable or
	// invalid.
	Play(*Player, []byte) (bool, error)

	// Retrieve the winning player.
	Winner() Player

	// Determine whether the match is a stalemate.
	Stalemate() bool

	// Advance the turn.
	NextPlayer() *Player
}
```
