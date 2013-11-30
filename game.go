package arena

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const URL = "http://localhost:5000/fourup"
const BaseUri = "https://battleofbits.com"

const Empty = 0
const Red = 1
const Black = 2

type FourUpMatch struct {
	Id          int64
	RedPlayer   *Player
	BlackPlayer *Player
	// Whose turn is it
	CurrentPlayer *Player
	Winner        int64
	Board         *[NumRows][NumColumns]int
}

type TurnPlayers struct {
	Red   string `json:"R"`
	Black string `json:"B"`
}

type FourUpTurn struct {
	Href     string                      `json:"href"`
	Players  *TurnPlayers                `json:"players"`
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

func serializeTurn(match *FourUpMatch) *FourUpTurn {
	return &FourUpTurn{
		Href:  getMatchHref(match.Id),
		Board: GetStringBoard(match.Board),
		Turn:  fmt.Sprintf(BaseUri+"/players/%s", match.CurrentPlayer.Name),
		Players: &TurnPlayers{
			Red:   fmt.Sprintf(BaseUri+"/players/%s", match.RedPlayer.Name),
			Black: fmt.Sprintf(BaseUri+"/players/%s", match.BlackPlayer.Name),
		},
	}
}

func CreateFourUpMatch(redPlayer *Player, blackPlayer *Player) *FourUpMatch {
	board := InitializeBoard()
	match := &FourUpMatch{
		RedPlayer:   redPlayer,
		BlackPlayer: blackPlayer,
		Board:       board,
		// Red plays first, I believe.
		CurrentPlayer: redPlayer,
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

func DoForfeit(loser *Player, reason error) {
	fmt.Println(fmt.Sprintf("player %s forfeits because of %s", loser.Username, reason.Error()))
}

func DoGameOver(match *FourUpMatch, winner *Player, loser *Player) {
	fmt.Println("Game is over. Winner is ", winner.Username, ". Notifying winner and loser...")
	MarkWinner(match, winner)
	NotifyWinner(winner)
	NotifyLoser(loser)
}

func getMatchHref(matchId int64) string {
	return fmt.Sprintf(BaseUri+"/games/four-up/matches/%d", matchId)
}

// Assemble and make an HTTP request to the user's URL
// Returns the column of the response
func GetMove(match *FourUpMatch) (int, error) {
	turn := serializeTurn(match)
	postBody, err := json.Marshal(turn)
	checkError(err)
	httpResponse, err := MakeRequest(match.CurrentPlayer.Url, postBody)
	return ParseResponse(httpResponse)
}

func MakeRequest(url string, postBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(postBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "battleofbits/0.1")
	// XXX, set a timeout here
	return client.Do(req)
}

// Retrieves the column from the http response
func ParseResponse(response *http.Response) (int, error) {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}
	var fourUpResponse FourUpResponse
	err = json.Unmarshal(body, &fourUpResponse)
	if err != nil {
		return -1, err
	}
	return fourUpResponse.Column, nil
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
	move, err := GetMove(match)
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
		match.CurrentPlayer = redPlayer
		err := DoPlayerMove(redPlayer, blackPlayer, match, 1)
		if err != nil {
			break
		}
		match.CurrentPlayer = blackPlayer
		err = DoPlayerMove(blackPlayer, redPlayer, match, 2)
		if err != nil {
			break
		}
	}
	return match
}
