package main

import (
	"fmt"
	"github.com/lib/pq"
)

func CreatePlayer(username string, name string) (int64, *pq.Error) {
	db := getConnection()
	result, err := db.Exec("INSERT INTO players (username, name) VALUES ($1, $2) RETURNING id", username, name)
	fmt.Println(result)
	fmt.Println("after insert")
	var pqerr *pq.Error
	if err != nil {
		pqerr = err.(*pq.Error)
	}
	if pqerr != nil && pqerr.Code.Name() == "unique_violation" {
		return -1, pqerr
	}
	checkError(err)
	id, err := result.LastInsertId()
	fmt.Println("checking insert id")
	checkError(err)
	return id, nil
}

func CreateGame() {
	playerOneId, err := CreatePlayer("Kevin Burke", "kevinburke")
	if err != nil && err.Code.Name() != "unique_violation" {
		checkError(err)
	}
	playerTwoId, err := CreatePlayer("Kyle Conroy", "kyleconroy")
	if err != nil && err.Code.Name() != "unique_violation" {
		checkError(err)
	}
	fmt.Println(playerOneId)
	fmt.Println(playerTwoId)
}

const URL = "http://localhost:5000"

func main() {
	CreateGame()
	fmt.Println("done")
}
