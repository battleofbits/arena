package fourup

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/battleofbits/arena/arena"
	"time"
)

// This level of indirection necessary to translate between int/string
// representation. Maybe we should just store everything as strings.
type FourUpBoard struct {
	Board [NumRows][NumColumns]int8
}

type FourUpMatch struct {
	Players       []*arena.Player
	Started       time.Time
	Board         *FourUpBoard
	currentPlayer *arena.Player
	MoveId        int64
	winner        *arena.Player
	Id            int64
}

func getStringBoard(board [NumRows][NumColumns]int8) [NumRows][NumColumns]string {
	var stringBoard [NumRows][NumColumns]string
	for row := int8(0); row < NumRows; row++ {
		for column := int8(0); column < NumColumns; column++ {
			if board[row][column] == Empty {
				stringBoard[row][column] = ""
			} else if board[row][column] == Red {
				stringBoard[row][column] = "R"
			} else if board[row][column] == Black {
				stringBoard[row][column] = "B"
			} else {
				panic(fmt.Sprint("invalid value ", board[row][column], " for a board"))
			}
		}
	}
	return stringBoard
}

// Leaving this here till we're sure we don't need it, the method below
// replaces this one.
//func serializeTurn(match *FourUpMatch) *FourUpTurn {
//return &FourUpTurn{
//Href:  getMatchHref(match.Id),
//Board: GetStringBoard(match.Board),
//Turn:  fmt.Sprintf(BaseUri+"/players/%s", match.CurrentPlayer.Name),
//Players: &TurnPlayers{
//Red:   fmt.Sprintf(BaseUri+"/players/%s", match.RedPlayer.Name),
//Black: fmt.Sprintf(BaseUri+"/players/%s", match.BlackPlayer.Name),
//},
//}
//}

func (b *FourUpBoard) MarshalJSON() ([]byte, error) {
	sbd := getStringBoard(b.Board)
	return json.Marshal(sbd)
}

// Retrieve the current player.
func (m *FourUpMatch) CurrentPlayer() *arena.Player {
	return m.currentPlayer
}

func (m *FourUpMatch) Stalemate() bool {
	// XXX
	return false
}

func CreateMatch(players []*arena.Player) (FourUpMatch, error) {
	if len(players) != 2 {
		msg := fmt.Sprintf("wrong number of players: %d", len(players))
		return FourUpMatch{}, errors.New(msg)
	}

	return createFourUpMatch(players), nil
}

// Apply the move to the board, write it to the database
// Returns a boolean (whether the game is over) and an error (whether the move
// was invalid)
func (m *FourUpMatch) Play(player *arena.Player, data []byte) (bool, error) {
	var fm fourUpMove
	err := json.Unmarshal(data, fm)
	if err != nil {
		return true, err
	}
	err = doNewMove(fm.Column, m)
	if err != nil {
		// XXX, assign the winner to be the other player.
		m.winner = player
		return true, err
	}
	if over, winner := gameOver(m.Board.Board); over {
		m.winner = player
		return true, nil
	}
	return false, nil
}

func createFourUpMatch(players []*arena.Player) *arena.Match {
	var board *FourUpBoard
	return &FourUpMatch{
		Players: players,
		Board:   board,
		// Red plays first, I believe.
		currentPlayer: players[0],
		MoveId:        0,
		Started:       time.Now().UTC(),
	}
}

func writeMatch(match *FourUpMatch) error {
	db := arena.GetConnection()
	defer db.Close()
	jsonBoard, err := json.Marshal(match.Board)
	if err != nil {
		return err
	}
	query := "INSERT INTO fourup_matches " +
		"(player_red, player_black, board, started) VALUES " +
		"($1, $2, $3, NOW() at time zone 'utc') RETURNING id"
	return db.QueryRow(query, match.Players[0].Id, match.Players[1].Id,
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

// Convert players => board color
func (m *FourUpMatch) getCurrentTurnColor() int8 {
	if m.CurrentPlayer() == m.Players[0] {
		return Red
	} else {
		return Black
	}
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
	return nil
}

// Update the match in the database
// Assumes the match has been initialized at some point
func UpdateMatch(match *FourUpMatch) error {
	db := arena.GetConnection()
	defer db.Close()
	jsonBoard, err := json.Marshal(match.Board)
	if err != nil {
		return err
	}
	query := "UPDATE fourup_matches SET board = $1 WHERE id = $2"
	_, err = db.Exec(query, string(jsonBoard), match.Id)
	return err
}

func markWinner(match *Match, winner *arena.Player) error {
	db := GetConnection()
	defer db.Close()
	_, err := db.Exec("UPDATE fourup_matches SET winner = $1, "+
		"finished = NOW() at time zone 'utc' WHERE id = $2",
		winner.Id, match.Id)
	return err
}
