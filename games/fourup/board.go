package fourup

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Constant representing an empty board section
const Empty = 0

// Constant representing red chit on the board
const Red = 1

// Constant representing black chit on the board
const Black = 2

// Number of rows in a fourup board
const NumRows = 6

// Number of columns in a fourup board
const NumColumns = 7

// Number of consecutive chits you need to win the game
const NumConsecutive = 4

// Checks fourup board state
// Author: Kevin Burke <kev@inburke.com>

type Board [NumRows][NumColumns]int8
type StringBoard [NumRows][NumColumns]string

func (b Board) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.getStringBoard())
}

func (b Board) isFull() bool {
	// will check the top row, which is always the last to fill up.
	for column := 0; column < NumColumns; column++ {
		if b[0][column] == Empty {
			return false
		}
	}
	return true
}

func (board Board) getStringBoard() StringBoard {
	var stringBoard StringBoard
	for row := int8(0); row < NumRows; row++ {
		for column := int8(0); column < NumColumns; column++ {
			if board[row][column] == Empty {
				stringBoard[row][column] = ""
			} else if board[row][column] == Red {
				stringBoard[row][column] = "R"
			} else if board[row][column] == Black {
				stringBoard[row][column] = "B"
			} else {
				panic(fmt.Sprint("invalid value ", board[row][column], " for a board"))
			}
		}
	}
	return stringBoard
}

// row varies, column does not.
func checkVerticalWin(column int8, board Board) (bool, int8) {
	checkRowInColumn := func(column int8, row int8,
		board [NumRows][NumColumns]int8) (bool, int8) {
		initColor := board[row][column]
		for k := int8(0); k < NumConsecutive; k++ {
			if row+k >= NumRows {
				return false, -1
			}
			value := board[row+k][column]
			if value == Empty || value != initColor {
				return false, -1
			}
		}
		// if we get here and haven't broken, seen 4 in a row of the same color
		return true, initColor
	}

	for row := int8(0); row <= (NumRows - NumConsecutive); row++ {
		initColor := board[row][column]
		if initColor == Empty {
			continue
		}
		if over, winner := checkRowInColumn(column, row, board); over {
			return true, winner
		}
	}
	return false, -1
}

func checkHorizontalWin(row int8, board Board) (bool, int8) {
	checkColumnInRow := func(row int8, column int8, board Board) (bool, int8) {
		initColor := board[row][column]
		for k := int8(0); k < NumConsecutive; k++ {
			if column+k >= NumColumns {
				return false, -1
			}
			if board[row][column+k] != initColor {
				return false, -1
			}
		}
		// if we get here and haven't broken, seen 4 in a row of the same color
		return true, initColor
	}
	for column := int8(0); column < NumConsecutive; column++ {
		initColor := board[row][column]
		if initColor == Empty {
			continue
		}
		if over, winner := checkColumnInRow(row, column, board); over {
			return true, winner
		}
	}
	return false, -1
}

// check squares down and to the right for a match
func checkSoutheastDiagonalWin(row int8, column int8, board Board) (bool, int8) {
	initColor := board[row][column]

	if initColor == Empty {
		return false, -1
	}
	for i := int8(0); i < NumConsecutive; i++ {
		if board[row+i][column+i] != initColor {
			return false, -1
		}
	}
	return true, initColor
}

func checkSouthwestDiagonalWin(row int8, column int8, board Board) (bool, int8) {
	initColor := board[row][column]
	if initColor == Empty {
		return false, -1
	}
	for i := int8(0); i < NumConsecutive; i++ {
		if board[row+i][column-i] != initColor {
			return false, -1
		}
	}
	return true, initColor
}

// Checks if a connect four exists
// I'm sure there's some more efficient way to conduct these checks, but at
// modern computer speeds, it really doesn't matter
func (board Board) gameOver() (bool, int8) {
	for column := int8(0); column < NumColumns; column++ {
		if over, winner := checkVerticalWin(column, board); over {
			return true, winner
		}
	}

	for row := int8(0); row < NumRows; row++ {
		if over, winner := checkHorizontalWin(row, board); over {
			return true, winner
		}
	}
	for row := int8(0); row <= (NumRows - NumConsecutive); row++ {
		for column := int8(0); column <= (NumColumns - NumConsecutive); column++ {
			if over, winner := checkSoutheastDiagonalWin(row, column, board); over {
				return true, winner
			}
		}
	}
	for column := int8(NumColumns - NumConsecutive); column < NumColumns; column++ {
		for row := int8(0); row <= (NumRows - NumConsecutive); row++ {
			if over, winner := checkSouthwestDiagonalWin(row, column, board); over {
				return true, winner
			}
		}
	}
	return false, -1
}

// Returns error if the move is invalid
func (bp *Board) applyMove(move int8, color int8) error {
	if move >= NumColumns || move < 0 {
		return errors.New(fmt.Sprintf("Move %d is invalid", move))
	}
	for i := NumRows - 1; i >= 0; i-- {
		if bp[i][move] == 0 {
			bp[i][move] = color
			return nil
		}
	}
	return errors.New(fmt.Sprintf("No room in column %d for a move", move))
}
