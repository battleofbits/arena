// Data layer for a match
package arena

import (
	"encoding/json"
	"fmt"
)

type FourUpMatch struct {
	Id          int64
	RedPlayer   *Player
	BlackPlayer *Player
	// Whose turn is it
	CurrentPlayer *Player
	Winner        int64
	Board         *[NumRows][NumColumns]int
	MoveId        int
}

func getMatchHref(matchId int64) string {
	return fmt.Sprintf(BaseUri+"/games/four-up/matches/%d", matchId)
}

func CreateFourUpMatch(redPlayer *Player, blackPlayer *Player) *FourUpMatch {
	board := InitializeBoard()
	match := &FourUpMatch{
		RedPlayer:   redPlayer,
		BlackPlayer: blackPlayer,
		Board:       board,
		// Red plays first, I believe.
		CurrentPlayer: redPlayer,
		MoveId:        0,
	}
	return match
}

func WriteMatch(match *FourUpMatch) error {
	db := getConnection()
	defer db.Close()
	stringBoard := GetStringBoard(match.Board)
	jsonBoard, err := json.Marshal(stringBoard)
	if err != nil {
		return err
	}
	query := "INSERT INTO fourup_matches " +
		"(player_red, player_black, board, started) VALUES " +
		"($1, $2, $3, NOW() at time zone 'utc') RETURNING id"
	return db.QueryRow(query, match.RedPlayer.Id, match.BlackPlayer.Id, string(jsonBoard)).Scan(&match.Id)
}

// Update the match in the database
// Assumes the match has been initialized at some point
func UpdateMatch(match *FourUpMatch) error {
	db := getConnection()
	defer db.Close()
	stringBoard := GetStringBoard(match.Board)
	jsonBoard, err := json.Marshal(stringBoard)
	if err != nil {
		return err
	}
	query := "UPDATE fourup_matches SET board = $1 WHERE id = $2"
	_, err = db.Exec(query, string(jsonBoard), match.Id)
	return err
}
