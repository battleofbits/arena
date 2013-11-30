package main

import (
	"fmt"
	"github.com/lib/pq"
	"math/rand"
)

func CreatePlayer(username string, name string) (*Player, error) {
	db := getConnection()
	_, err := db.Exec("INSERT INTO players (username, name) VALUES ($1, $2) RETURNING id", username, name)
	var pqerr *pq.Error
	if err != nil {
		pqerr = err.(*pq.Error)
	}
	if pqerr != nil && pqerr.Code.Name() == "unique_violation" {
		return &Player{}, pqerr
	}
	checkError(err)
	// XXX
	return GetPlayerByName(name)
}

type Player struct {
	// The autoid for the player
	Id int64
	// The player's unique Id
	Name string
	// The player's friendly name
	Username string
	Url      string
}

func GetPlayerByName(name string) (*Player, error) {
	var p Player
	db := getConnection()
	err := db.QueryRow("SELECT * FROM players WHERE name = $1", name).Scan(&p.Id, &p.Username, &p.Name, &p.Url)
	if err != nil {
		return &Player{}, err
	} else {
		return &p, nil
	}
}

type FourUpMatch struct {
	RedPlayerId   int64
	BlackPlayerId int64
	Winner        int64
	Board         *[6][7]int
}

func CreateFourUpMatch(redPlayer *Player, blackPlayer *Player) (*FourUpMatch, error) {
	board := InitializeBoard()
	db := getConnection()
	_, err := db.Exec("INSERT INTO fourup_matches (player_red, player_black) VALUES ($1, $2) RETURNING id", redPlayer.Id, blackPlayer.Id)
	checkError(err)
	return &FourUpMatch{
		RedPlayerId:   redPlayer.Id,
		BlackPlayerId: blackPlayer.Id,
		Board:         board,
	}, nil
}

func DoForfeit(loser *Player, reason error) {
	fmt.Println(fmt.Sprintf("player %s forfeits because of %s", loser.Username, reason.Error()))
}

func DoGameOver(match *FourUpMatch, winner *Player, loser *Player) {
	MarkWinner(match, winner)
	NotifyWinner(winner)
	NotifyLoser(loser)
}

func GetMove(player *Player) (int, error) {
	return rand.Intn(7), nil
}

func NotifyWinner(winner *Player) {

}

func NotifyLoser(loser *Player) {

}

func MarkWinner(match *FourUpMatch, winner *Player) {

}

func DoPlayerMove(player *Player, otherPlayer *Player, match *FourUpMatch, playerId int) error {
	move, err := GetMove(player)
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
		return err
	}
	if IsBoardFull(*match.Board) {
		DoTieGame(match, player, otherPlayer)
		return err
	}
	return nil
}

func DoTieGame(match *FourUpMatch, playerOne *Player, playerTwo *Player) {

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

const URL = "http://localhost:5000"

func main() {
	redPlayer, _ := CreatePlayer("Kevin Burke", "kevinburke")
	blackPlayer, _ := CreatePlayer("Kyle Conroy", "kyleconroy")
	match, fourupErr := CreateFourUpMatch(redPlayer, blackPlayer)
	checkError(fourupErr)
	match = DoMatch(match, redPlayer, blackPlayer)
	fmt.Println(match.Board)
	fmt.Println("done")
}
