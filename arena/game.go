package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"io/ioutil"
	"net/http"
)

type Player struct {
	// The autoid for the player
	Id int64
	// The player's unique Id
	Name string
	// The player's friendly name
	Username string
	Url      string
}

type FourUpMatch struct {
	Id            int64
	RedPlayerId   int64
	BlackPlayerId int64
	Winner        int64
	Board         *[NumRows][NumColumns]int
}

type FourUpTurn struct {
	Href     string                      `json:"href"`
	Players  map[string]string           `json:"players"`
	Turn     string                      `json:"turn"`
	Loser    string                      `json:"loser"`
	Winner   string                      `json:"winner"`
	Started  string                      `json:"started"`
	Finished string                      `json:"finished"`
	Moves    string                      `json:"moves"`
	Board    [NumRows][NumColumns]string `json:"board"`
}

type FourUpResponse struct {
	Column int `json:"column"`
}

const URL = "http://localhost:5000/fourup"

const Empty = 0
const Red = 1
const Black = 2

func CreatePlayer(username string, name string, url string) (*Player, error) {
	db := getConnection()
	defer db.Close()
	player := &Player{
		Username: username,
		Name:     name,
		Url:      url,
	}
	err := db.QueryRow("INSERT INTO players (username, name, url) VALUES ($1, $2, $3) RETURNING id", username, name, url).Scan(&player.Id)
	var pqerr *pq.Error
	if err != nil {
		pqerr = err.(*pq.Error)
	}
	if pqerr != nil && pqerr.Code.Name() == "unique_violation" {
		return &Player{}, pqerr
	}
	checkError(err)
	return player, nil
}

func GetPlayerByName(name string) (*Player, error) {
	var p Player
	db := getConnection()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM players WHERE name = $1", name).Scan(&p.Id, &p.Username, &p.Name, &p.Url)
	if err != nil {
		return &Player{}, err
	} else {
		return &p, nil
	}
}

func CreateFourUpMatch(redPlayer *Player, blackPlayer *Player) (*FourUpMatch, error) {
	board := InitializeBoard()
	match := &FourUpMatch{
		RedPlayerId:   redPlayer.Id,
		BlackPlayerId: blackPlayer.Id,
		Board:         board,
	}
	db := getConnection()
	defer db.Close()
	err := db.QueryRow("INSERT INTO fourup_matches "+
		"(player_red, player_black, started) VALUES "+
		"($1, $2, NOW() at time zone 'utc') RETURNING id",
		redPlayer.Id, blackPlayer.Id).Scan(&match.Id)
	checkError(err)
	return match, nil
}

func DoForfeit(loser *Player, reason error) {
	fmt.Println(fmt.Sprintf("player %s forfeits because of %s", loser.Username, reason.Error()))
}

func DoGameOver(match *FourUpMatch, winner *Player, loser *Player) {
	fmt.Println("Game is over. Winner is ", winner.Username, ". Notifying winner and loser...")
	MarkWinner(match, winner)
	NotifyWinner(winner)
	NotifyLoser(loser)
}

func getHref(id int64) string {
	return fmt.Sprintf("https://battleofbits.com/games/four-up/matches/%d", id)
}

func serializeTurn(match *FourUpMatch) *FourUpTurn {
	return &FourUpTurn{
		Href:  getHref(match.Id),
		Board: GetStringBoard(match.Board),
	}
}

func GetMove(player *Player, match *FourUpMatch) (int, error) {
	turn := serializeTurn(match)
	postBody, err := json.Marshal(turn)
	checkError(err)
	req, err := http.NewRequest("POST", player.Url, bytes.NewReader(postBody))
	if err != nil {
		return -1, err
	}
	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "battleofbits/0.1")
	httpResponse, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return -1, err
	}
	var response FourUpResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}
	return response.Column, nil
}

func NotifyWinner(winner *Player) {
	fmt.Println("Notifying winner...")
}

func NotifyLoser(loser *Player) {
	fmt.Println("Notifying loser...")
}

func MarkWinner(match *FourUpMatch, winner *Player) error {
	db := getConnection()
	defer db.Close()
	_, err := db.Exec("UPDATE fourup_matches SET winner = $1, "+
		"finished = NOW() at time zone 'utc' WHERE id = $2",
		winner.Id, match.Id)
	return err
}

func DoPlayerMove(player *Player, otherPlayer *Player, match *FourUpMatch, playerId int) error {
	move, err := GetMove(player, match)
	if err != nil {
		DoForfeit(player, err)
		DoGameOver(match, otherPlayer, player)
		return err
	}
	match.Board, err = ApplyMoveToBoard(move, playerId, match.Board)
	if err != nil {
		DoForfeit(player, err)
		DoGameOver(match, otherPlayer, player)
		return err
	}
	if GameOver(*match.Board) {
		DoGameOver(match, player, otherPlayer)
		return errors.New("Game is over.")
	}
	if IsBoardFull(*match.Board) {
		DoTieGame(match, player, otherPlayer)
		return err
	}
	return nil
}

func DoTieGame(match *FourUpMatch, playerOne *Player, playerTwo *Player) {
	fmt.Println("Tie Game!")
}

func DoMatch(match *FourUpMatch, redPlayer *Player, blackPlayer *Player) *FourUpMatch {
	for {
		err := DoPlayerMove(redPlayer, blackPlayer, match, 1)
		if err != nil {
			break
		}
		err = DoPlayerMove(blackPlayer, redPlayer, match, 2)
		if err != nil {
			break
		}
	}
	return match
}

func main() {
	redPlayer, _ := CreatePlayer("Kevin Burke", "kevinburke", URL)
	blackPlayer, _ := CreatePlayer("Kyle Conroy", "kyleconroy", URL)
	match, fourupErr := CreateFourUpMatch(redPlayer, blackPlayer)
	checkError(fourupErr)
	match = DoMatch(match, redPlayer, blackPlayer)
	fmt.Println(match.Board)
	fmt.Println("done")
}
