package main

import (
	"fmt"
	"github.com/lib/pq"
)

func CreatePlayer(username string, name string) (int64, *pq.Error) {
	db := getConnection()
	var id int64
	_, err := db.Exec("INSERT INTO players (username, name) VALUES ($1, $2) RETURNING id", username, name)
	fmt.Println("after insert")
	var pqerr *pq.Error
	if err != nil {
		pqerr = err.(*pq.Error)
	}
	if pqerr != nil && pqerr.Code.Name() == "unique_violation" {
		return -1, pqerr
	}
	checkError(err)
	return id, nil
}

type Player struct {
	// The player's unique Id
	Name string
	// The player's friendly name
	Username string
	// The autoid for the player
	Id int64
}

func GetPlayerByName(name string) (*Player, error) {
	var p Player
	db := getConnection()
	err := db.QueryRow("SELECT * FROM players WHERE name = $1", name).Scan(&p.Username, &p.Name, &p.Id)
	fmt.Println(p.Username)
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
	Board         [7][7]int
}

func InitializeBoard() [7][7]int {
	var board [7][7]int
	// Board is initialized to be filled with zeros.
	return board
}

func CreateFourUpMatch(redPlayerName string, blackPlayerName string) (*FourUpMatch, error) {
	redPlayer, err := GetPlayerByName(redPlayerName)
	if err != nil {
		return &FourUpMatch{}, err
	}
	blackPlayer, err := GetPlayerByName(blackPlayerName)
	if err != nil {
		return &FourUpMatch{}, err
	}
	board := InitializeBoard()
	db := getConnection()
	_, err = db.Exec("INSERT INTO fourup_matches (player_red, player_black) VALUES ($1, $2) RETURNING id", redPlayer.Id, blackPlayer.Id)
	checkError(err)
	fmt.Println("returning match")
	return &FourUpMatch{
		RedPlayerId:   redPlayer.Id,
		BlackPlayerId: blackPlayer.Id,
		Board:         board,
	}, nil
}

const URL = "http://localhost:5000"

func main() {
	_, err := CreatePlayer("Kevin Burke", "kevinburke")
	if err != nil && err.Code.Name() != "unique_violation" {
		checkError(err)
	}
	_, err = CreatePlayer("Kyle Conroy", "kyleconroy")
	if err != nil && err.Code.Name() != "unique_violation" {
		checkError(err)
	}
	match, fourupErr := CreateFourUpMatch("kevinburke", "kyleconroy")
	checkError(fourupErr)
	fmt.Println(match.Board)
	fmt.Println("done")
}
