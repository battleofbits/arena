// Data layer for a match
package arena

import (
	"encoding/json"
	"fmt"
	"time"
)

// Exactly the same interface as sql.NullString but with a MarshalJSON method
type NullString struct {
	String string
	Valid  bool
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte{}, nil
	} else {
		return json.Marshal(n.String)
	}
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte{}, nil
	} else {
		return json.Marshal(n.Time)
	}
}

type FourUpMatch struct {
	Id          int64
	Started     time.Time
	Finished    time.Time
	RedPlayer   *Player
	BlackPlayer *Player
	// Whose turn is it
	CurrentPlayer *Player
	Winner        *Player
	Board         *[NumRows][NumColumns]int8
	MoveId        int
}

// public facing thingy
type MatchResponse struct {
	Id          int64                      `json:"id"`
	CurrentMove string                     `json:"current_move"`
	Winner      *NullString                `json:"winner"`
	RedPlayer   string                     `json:"red_player"`
	BlackPlayer string                     `json:"black_player"`
	Board       *[NumRows][NumColumns]int8 `json:"board"`
	Started     *NullTime                  `json:"started"`
	Finished    *NullTime                  `json:"finished"`
}

func (m *FourUpMatch) MarshalJSON() ([]byte, error) {
	winnerString := &NullString{
		Valid: true,
	}
	if m.Winner != nil {
		if m.Winner == m.RedPlayer {
			winnerString.String = m.RedPlayer.Name
		} else {
			winnerString.String = m.BlackPlayer.Name
		}
	}
	startNullable := &NullTime{
		Valid: true,
		Time:  m.Started,
	}
	finishedNullable := &NullTime{
		Valid: true,
		Time:  m.Finished,
	}
	var currentPlayerName, redPlayerName, blackPlayerName string
	if m.CurrentPlayer != nil {
		currentPlayerName = m.CurrentPlayer.Name
	}
	if m.RedPlayer != nil {
		redPlayerName = m.RedPlayer.Name
	}
	if m.BlackPlayer != nil {
		blackPlayerName = m.BlackPlayer.Name
	}
	return json.Marshal(&MatchResponse{
		Id:          m.Id,
		CurrentMove: currentPlayerName,
		Winner:      winnerString,
		Started:     startNullable,
		Finished:    finishedNullable,
		Board:       m.Board,
		RedPlayer:   redPlayerName,
		BlackPlayer: blackPlayerName,
	})
}

func (m *FourUpMatch) GetCurrentTurnColor() int8 {
	if m.CurrentPlayer == m.RedPlayer {
		return Red
	} else {
		return Black
	}
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
		Started:       time.Now().UTC(),
	}
	return match
}

func WriteMatch(match *FourUpMatch) error {
	db := GetConnection()
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
	db := GetConnection()
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

func GetMatches() ([]*FourUpMatch, error) {
	db := GetConnection()
	defer db.Close()
	query := "select red_players.name as red_player, " +
		"black_players.name as black_player, " +
		"winners.name as winner, " +
		"fourup_matches.id, fourup_matches.started, fourup_matches.finished " +
		"from fourup_matches " +
		"inner join players as red_players on red_players.id=fourup_matches.player_red " +
		"inner join players as black_players on black_players.id=fourup_matches.player_black " +
		"inner join players as winners on winners.id=fourup_matches.winner order by started limit 50"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var matches []*FourUpMatch
	for rows.Next() {
		var m FourUpMatch
		var redName string
		var blackName string
		var winnerName string
		err = rows.Scan(&redName, &blackName, &winnerName, &m.Id, &m.Started,
			&m.Finished)
		if err != nil {
			return nil, err
		}
		m.RedPlayer = &Player{
			Name: redName,
		}
		m.BlackPlayer = &Player{
			Name: blackName,
		}
		if winnerName == redName {
			m.Winner = m.RedPlayer
		} else if winnerName == blackName {
			m.Winner = m.BlackPlayer
		}
		matches = append(matches, &m)
	}
	return matches, nil
}
