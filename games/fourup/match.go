package fourup

import (
	"bytes"
	"github.com/battleofbits/arena/arena"
)

func CreateMatch(players []*arena.Player) (*arena.Match, error) {
	if len(players) != 2 {
		return nil, errors.New("wrong number of players: %d", len(players))
	}
	match := createFourUpMatch(players)
	dbErr := writeMatch(match)
	if dbErr != nil {
		return dbErr
	}
	return nil
}

// Apply the move to the board, write it to the database
// Returns a boolean (whether the game is over) and an error (whether the move
// was invalid)
func (m *Match) Play(player *arena.Player, data []bytes) (bool, error) {
	var fm fourUpMove
	err := json.Unmarshal(data, fm)
	if err != nil {
		return true, err
	}
	err = doNewMove(fm.Column, m)
	if err != nil {
		// XXX, assign the winner to be the other player.
		match.Winner = player - 1
		return true, err
	}
	if over, winner := gameOver(*match.Board); over {
		match.Winner = player
		return true, nil
	}
	return false, nil
}

func createFourUpMatch(players []*arena.Player) *arena.Match {
	board := initializeBoard()
	return &arena.Match{
		Players: players,
		Board:   board,
		// Red plays first, I believe.
		CurrentPlayer: players[0],
		MoveId:        0,
		Started:       time.Now().UTC(),
	}
}

func writeMatch(match *FourUpMatch) error {
	db := arena.GetConnection()
	defer db.Close()
	stringBoard := GetStringBoard(match.Board)
	jsonBoard, err := json.Marshal(stringBoard)
	if err != nil {
		return err
	}
	query := "INSERT INTO fourup_matches " +
		"(player_red, player_black, board, started) VALUES " +
		"($1, $2, $3, NOW() at time zone 'utc') RETURNING id"
	return db.QueryRow(query, match.RedPlayer.Id, match.BlackPlayer.Id,
		string(jsonBoard)).Scan(&match.Id)
}

var moveWriter = writeMove

// Write a new move to the database
func writeMove(move int8, match *arena.Match) (int64, error) {
	db := arena.GetConnection()
	defer db.Close()
	query := "INSERT INTO fourup_moves (fourup_column, player, move_number, match_id, played)" +
		"VALUES ($1, $2, $3, $4, NOW() at time zone 'utc') RETURNING id"
	var moveId int64
	err := db.QueryRow(query, move, match.CurrentPlayer.Id, match.MoveId, match.Id).Scan(&moveId)
	return moveId, err
}

// A four up move
type fourUpMove struct {
	Column int8 `json:"column"`
}

// Do a whole bunch of stuff associated with new moves
// Error handling is a little tricky because most of the errors would be
// database or other errors.
func doNewMove(column int8, match *arena.Match) error {
	var err error
	match.Board, err = applyMoveToBoard(column, match.GetCurrentTurnColor(),
		match.Board)
	// XXX
	//if err != nil {
	//DoForfeit(player, err)
	//DoGameOver(match, otherPlayer, player)
	//return err
	//}
	// once we know move was valid, update the database
	_, err = writeMove(column, match)
	if err != nil {
		return err
	}
	match.MoveId++
	err = UpdateMatch(match)
	if err != nil {
		return err
	}
}
