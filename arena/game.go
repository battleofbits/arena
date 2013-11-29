package main

import (
	"fmt"
	"github.com/lib/pq"
)

func CreatePlayer(username string, name string) (*Player, error) {
	db := getConnection()
	_, err := db.Exec("INSERT INTO players (username, name) VALUES ($1, $2) RETURNING id", username, name)
	fmt.Println("after insert")
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

func CreateFourUpMatch(redPlayer *Player, blackPlayer *Player) (*FourUpMatch, error) {
	board := InitializeBoard()
	db := getConnection()
	_, err := db.Exec("INSERT INTO fourup_matches (player_red, player_black) VALUES ($1, $2) RETURNING id", redPlayer.Id, blackPlayer.Id)
	checkError(err)
	fmt.Println("returning match")
	return &FourUpMatch{
		RedPlayerId:   redPlayer.Id,
		BlackPlayerId: blackPlayer.Id,
		Board:         board,
	}, nil
}

func DoForfeit(loser *Player, reason error) {

}

func DoGameOver(match *FourUpMatch, winner *Player, loser *Player) {
	MarkWinner(match, winner)
	NotifyWinner(winner)
	NotifyLoser(loser)
}

func GetMove(player *Player) (int, error) {
	return 3, nil
}

func NotifyWinner(winner *Player) {

}

func NotifyLoser(loser *Player) {

}

func ApplyMoveToBoard(move int, match *FourUpMatch) error {
	return nil
}

func MarkWinner(match *FourUpMatch, winner *Player) {

}

// row varies, column does not.
func checkVerticalWin(column int, board [7][7]int) bool {
	checkRowInColumn := func(column int, row int, board [7][7]int) bool {
		initColor := board[row][column]
		for k := 0; k < 4; k++ {
			if row+k >= 7 {
				return false
			}
			value := board[row+k][column]
			if value == 0 || value != initColor {
				return false
			}
		}
		// if we get here and haven't broken, seen 4 in a row of the same color
		return true
	}

	for row := 0; row < 4; row++ {
		initColor := board[row][column]
		if initColor == 0 {
			continue
		}
		if checkRowInColumn(column, row, board) {
			return true
		}
	}
	return false
}

func checkHorizontalWin(row int, board [7][7]int) bool {
	checkColumnInRow := func(row int, column int, board [7][7]int) bool {
		initColor := board[row][column]
		for k := 0; k < 4; k++ {
			if column+k >= 7 {
				return false
			}
			if board[row][column+k] != initColor {
				return false
			}
		}
		// if we get here and haven't broken, seen 4 in a row of the same color
		return true
	}
	for column := 0; column < 4; column++ {
		initColor := board[row][column]
		if initColor == 0 {
			continue
		}
		if checkColumnInRow(row, column, board) {
			return true
		}
	}
	return false
}

// check squares down and to the right for a match
func checkSoutheastDiagonalWin(row int, column int, board [7][7]int) bool {
	initColor := board[row][column]
	if initColor == 0 {
		return false
	}
	for i := 0; i < 4; i++ {
		if board[row+i][column+i] != initColor {
			return false
		}
	}
	return true
}

func checkSouthwestDiagonalWin(row int, column int, board [7][7]int) bool {
	initColor := board[row][column]
	if initColor == 0 {
		return false
	}
	for i := 0; i < 4; i++ {
		if board[row-i][column+i] != initColor {
			return false
		}
	}
	return true
}

// Checks if a connect four exists
// I'm sure there's some more efficient way to conduct these checks, but at
// modern computer speeds, it really doesn't matter
func GameOver(board [7][7]int) bool {
	for i := 0; i < 7; i++ {
		if checkVerticalWin(i, board) {
			return true
		}
		if checkHorizontalWin(i, board) {
			return true
		}
	}
	for row := 0; row < 4; row++ {
		for column := 0; column < 4; column++ {
			if checkSoutheastDiagonalWin(row, column, board) {
				return true
			}
		}
	}
	for row := 3; row < 7; row++ {
		for column := 0; column < 4; column++ {
			if checkSouthwestDiagonalWin(row, column, board) {
				return true
			}
		}
	}
	return false
}

func DoMatch(match *FourUpMatch, redPlayer *Player, blackPlayer *Player) {
	for {
		move, err := GetMove(redPlayer)
		if err != nil {
			DoForfeit(redPlayer, err)
			DoGameOver(match, blackPlayer, redPlayer)
		}
		err = ApplyMoveToBoard(move, match)
		if GameOver(match.Board) {
			DoGameOver(match, redPlayer, blackPlayer)
			break
		}
		move, err = GetMove(blackPlayer)
		if err != nil {
			DoForfeit(redPlayer, err)
			DoGameOver(match, blackPlayer, redPlayer)
		}
		ApplyMoveToBoard(move, match)
		err = ApplyMoveToBoard(move, match)
		if GameOver(match.Board) {
			DoGameOver(match, blackPlayer, redPlayer)
			break
		}
	}
}

const URL = "http://localhost:5000"

func main() {
	redPlayer, _ := CreatePlayer("Kevin Burke", "kevinburke")
	blackPlayer, _ := CreatePlayer("Kyle Conroy", "kyleconroy")
	match, fourupErr := CreateFourUpMatch(redPlayer, blackPlayer)
	checkError(fourupErr)
	DoMatch(match, redPlayer, blackPlayer)
	fmt.Println(match.Board)
	fmt.Println("done")
}
