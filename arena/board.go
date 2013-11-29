package main

import (
	"errors"
	"fmt"
)

// Checks fourup board state
// Author: Kevin Burke <kev@inburke.com>

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

func IsBoardFull(board [7][7]int) bool {
	// will check the top row, which is always the last to fill up.
	for column := 0; column < 7; column++ {
		if board[0][column] == 0 {
			return false
		}
	}
	return true
}

// Returns error if the move is invalid
func ApplyMoveToBoard(move int, playerId int, bp *[7][7]int) (*[7][7]int, error) {
	var badBoard *[7][7]int
	if move >= 7 || move < 0 {
		return badBoard, errors.New(fmt.Sprintf("Move %d is invalid", move))
	}
	for i := 6; i >= 0; i-- {
		if bp[i][move] == 0 {
			bp[i][move] = playerId
			return bp, nil
		}
	}
	return badBoard, errors.New(fmt.Sprintf("No room in column %d for a move", move))
}
