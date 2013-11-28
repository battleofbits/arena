package main

import (
	//"database/sql"
	"fmt"
)

func CreatePlayer(username string, name string) int64 {
	db := getConnection()
	result, err := db.Exec("INSERT INTO players (username, name) VALUES ($1, $2)", username, name)
	checkError(err)
	id, err := result.LastInsertId()
	checkError(err)
	return id
}

func CreateGame() {
	//db := getConnection()
	//err := db.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&u.Id, &u.Name, &u.Score)
	CreatePlayer("Kevin Burke", "kevinburke")
	CreatePlayer("Kyle Conroy", "kyleconroy")
}

const URL = "http://localhost:5000"

func main() {
	CreateGame()
	fmt.Println("done")
}
