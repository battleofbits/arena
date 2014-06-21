package fourup

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/battleofbits/arena/arena"
	"github.com/battleofbits/arena/engine"
	"time"
)

// This level of indirection necessary to translate between int/string
// representation. Maybe we should just store everything as strings.
type Match struct {
	Players       []*engine.Player
	Started       time.Time
	Board         *Board
	currentPlayer *engine.Player
	MoveId        int64
	winner        *engine.Player
	Id            int64
}

// A four up move
type Move struct {
	Column int8 `json:"column"`
}

// Retrieve the current player.
func (m *Match) CurrentPlayer() *engine.Player {
	return m.currentPlayer
}

func (m *Match) Stalemate() bool {
	// XXX
	return false
}

func (m *Match) Winner() engine.Player {
	// XXX
	return *m.winner
}

func (m *Match) NextPlayer() *engine.Player {
	if m.currentPlayer == m.Players[0] {
		m.currentPlayer = m.Players[1]
	} else {
		m.currentPlayer = m.Players[0]
	}
	return m.currentPlayer
}

func CreateMatch(players []*engine.Player) (*Match, error) {
	if len(players) != 2 {
		msg := fmt.Sprintf("wrong number of players: %d", len(players))
		return &Match{}, errors.New(msg)
	}

	board := &Board{}

	return &Match{
		Players: players,
		Board:   board,
		// Red plays first, I believe.
		currentPlayer: players[0],
		MoveId:        0,
		Started:       time.Now().UTC(),
	}, nil
}

// Serialize the move from the user's response and apply the move to the board.
// Returns a boolean (whether the game is over) and an error (whether the move
// was invalid).
func (m *Match) Play(player *engine.Player, data []byte) (bool, error) {
	var move Move

	err := json.Unmarshal(data, &move)
	if err != nil {
		return true, err
	}

	err = m.Board.applyMove(move.Column, m.getCurrentTurnColor())

	if err != nil {
		// XXX, assign the winner to be the other player.
		m.winner = player
		return true, err
	}

	if over, _ := m.Board.gameOver(); over {
		m.winner = player
		return true, nil
	}

	return false, nil
}

// Convert players => board color
func (m *Match) getCurrentTurnColor() int8 {
	if m.CurrentPlayer() == m.Players[0] {
		return Red
	} else {
		return Black
	}
}

// Update the match in the database.
// Assumes the match has been initialized at some point
func updateMatch(match *Match) error {
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
	db := arena.GetConnection()
	defer db.Close()
	_, err := db.Exec("UPDATE fourup_matches SET winner = $1, "+
		"finished = NOW() at time zone 'utc' WHERE id = $2",
		winner.Id, match.Id)
	return err
}
