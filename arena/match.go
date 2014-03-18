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

//func CreateMatch(*Player, *Player) *Match {

//}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	} else {
		return json.Marshal(n.String)
	}
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" {
		*n = NullString{
			Valid: false,
		}
		return nil
	} else {
		var out string
		err := json.Unmarshal(data, &out)
		if err != nil {
			return err
		}
		*n = NullString{
			String: out,
			Valid:  true,
		}
		return nil
	}
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	} else {
		return json.Marshal(n.Time)
	}
}

func (n *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*n = NullTime{
			Valid: false,
		}
		return nil
	} else {
		var out time.Time
		err := json.Unmarshal(data, &out)
		if err != nil {
			return err
		}
		*n = NullTime{
			Time:  out,
			Valid: true,
		}
		return nil
	}
}

// public facing thingy
//type MatchResponse struct {
//Id          int64                      `json:"id"`
//CurrentMove string                     `json:"current_move"`
//Winner      *NullString                `json:"winner"`
//RedPlayer   string                     `json:"red_player"`
//BlackPlayer string                     `json:"black_player"`
//Board       *[NumRows][NumColumns]int8 `json:"board"`
//Started     *NullTime                  `json:"started"`
//Finished    *NullTime                  `json:"finished"`
//Href        string                     `json:"href"`
//}

//func (mr *MatchResponse) SetHref() {
//mr.Href = fmt.Sprintf("https://battleofbits.com/games/four-up/matches/%d",
//mr.Id)
//}

//func (m *FourUpMatch) MarshalJSON() ([]byte, error) {
//winnerString := &NullString{
//Valid: true,
//}
//if m.Winner != nil {
//if m.Winner == m.RedPlayer {
//winnerString.String = m.RedPlayer.Name
//} else {
//winnerString.String = m.BlackPlayer.Name
//}
//}
//startNullable := &NullTime{
//Valid: true,
//Time:  m.Started,
//}
//finishedNullable := &NullTime{
//Valid: true,
//Time:  m.Finished,
//}
//var currentPlayerName, redPlayerName, blackPlayerName string
//if m.CurrentPlayer != nil {
//currentPlayerName = m.CurrentPlayer.Name
//}
//if m.RedPlayer != nil {
//redPlayerName = m.RedPlayer.Name
//}
//if m.BlackPlayer != nil {
//blackPlayerName = m.BlackPlayer.Name
//}
//mr := MatchResponse{
//Id:          m.Id,
//CurrentMove: currentPlayerName,
//Winner:      winnerString,
//Started:     startNullable,
//Finished:    finishedNullable,
//Board:       m.Board,
//RedPlayer:   redPlayerName,
//BlackPlayer: blackPlayerName,
//}
//mr.SetHref()
//return json.Marshal(&mr)
//}

func getMatchHref(matchId int64) string {
	return fmt.Sprintf(BaseUri+"/games/four-up/matches/%d", matchId)
}

// dry up the query a little bit
func getMatchQuery(singleId bool) string {
	var where string
	if singleId {
		where = "WHERE fourup_matches.id = $1 "
	} else {
		where = " "
	}
	return "select red_players.name as red_player, " +
		"black_players.name as black_player, " +
		"winners.name as winner, " +
		"fourup_matches.id, fourup_matches.started, fourup_matches.finished, " +
		"fourup_matches.board " +
		"from fourup_matches " +
		"inner join players as red_players " +
		"on red_players.id=fourup_matches.player_red " +
		"inner join players as black_players " +
		"on black_players.id=fourup_matches.player_black " +
		"inner join players as winners " +
		"on winners.id=fourup_matches.winner " +
		where + "order by started limit 50"
}

//func GetMatch(id int) (*FourUpMatch, error) {
//db := GetConnection()
//defer db.Close()
//query := getMatchQuery(true)
//fmt.Println(query)
//row := db.QueryRow(query, id)
//match, err := serializeMatch(row)
//if err != nil {
//return nil, err
//}
//return match, nil
//}

type MatchScanner interface {
	Scan(dest ...interface{}) error
}

//func serializeMatch(scanner MatchScanner) (*FourUpMatch, error) {
//var m FourUpMatch
//var redName string
//var blackName string
//var winnerName string
//var byteBoard []byte
//err := scanner.Scan(&redName, &blackName, &winnerName, &m.Id, &m.Started,
//&m.Finished, &byteBoard)
//if err != nil {
//return nil, err
//}
//board, err := GetIntBoard(byteBoard)
//if err != nil {
//return nil, err
//}
//m.Board = board
//m.RedPlayer = &Player{
//Name: redName,
//}
//m.BlackPlayer = &Player{
//Name: blackName,
//}
//if winnerName == redName {
//m.Winner = m.RedPlayer
//} else if winnerName == blackName {
//m.Winner = m.BlackPlayer
//}
//return &m, nil
//}

//func GetMatches() ([]*FourUpMatch, error) {
//db := GetConnection()
//defer db.Close()
//query := getMatchQuery(false)
//rows, err := db.Query(query)
//if err != nil {
//return nil, err
//}
//var matches []*FourUpMatch
//for rows.Next() {
//match, err := serializeMatch(rows)
//if err != nil {
//return nil, err
//}
//matches = append(matches, match)
//}
//return matches, nil
//}
