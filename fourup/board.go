package fourup

import (
	"errors"
	"fmt"
)

const Empty = 0
const Red = 1
const Black = 2

const NumRows = 6
const NumColumns = 7
const NumConsecutive = 4

// Checks fourup board state
type Board [NumRows][NumColumns]int

// row varies, column does not.
func (board *Board) checkVerticalWin(column int) bool {
	checkRowInColumn := func(column int, row int, board *Board) bool {
		initColor := board[row][column]
		for k := 0; k < NumConsecutive; k++ {
			if row+k >= NumRows {
				return false
			}
			value := board[row+k][column]
			if value == Empty || value != initColor {
				return false
			}
		}
		// if we get here and haven't broken, seen 4 in a row of the same color
		return true
	}

	for row := 0; row <= (NumRows - NumConsecutive); row++ {
		initColor := board[row][column]
		if initColor == Empty {
			continue
		}
		if checkRowInColumn(column, row, board) {
			return true
		}
	}
	return false
}

func (board *Board) checkHorizontalWin(row int) bool {
	checkColumnInRow := func(row int, column int, board *Board) bool {
		initColor := board[row][column]
		for k := 0; k < NumConsecutive; k++ {
			if column+k >= NumColumns {
				return false
			}
			if board[row][column+k] != initColor {
				return false
			}
		}
		// if we get here and haven't broken, seen 4 in a row of the same color
		return true
	}
	for column := 0; column < NumConsecutive; column++ {
		initColor := board[row][column]
		if initColor == Empty {
			continue
		}
		if checkColumnInRow(row, column, board) {
			return true
		}
	}
	return false
}

// check squares down and to the right for a match
func (board *Board) checkSoutheastDiagonalWin(row int, column int) bool {
	initColor := board[row][column]
	if initColor == Empty {
		return false
	}
	for i := 0; i < NumConsecutive; i++ {
		if board[row+i][column+i] != initColor {
			return false
		}
	}
	return true
}

func (board *Board) checkSouthwestDiagonalWin(row int, column int) bool {
	initColor := board[row][column]
	if initColor == Empty {
		return false
	}
	for i := 0; i < NumConsecutive; i++ {
		if board[row+i][column-i] != initColor {
			return false
		}
	}
	return true
}

// Checks if a connect four exists
// I'm sure there's some more efficient way to conduct these checks, but at
// modern computer speeds, it really doesn't matter
func (board *Board) GameOver() bool {
	for column := 0; column < NumColumns; column++ {
		if board.checkVerticalWin(column) {
			return true
		}
	}

	for row := 0; row < NumRows; row++ {
		if board.checkHorizontalWin(row) {
			return true
		}
	}
	for row := 0; row <= (NumRows - NumConsecutive); row++ {
		for column := 0; column <= (NumColumns - NumConsecutive); column++ {
			if board.checkSoutheastDiagonalWin(row, column) {
				return true
			}
		}
	}
	for column := (NumColumns - NumConsecutive); column < NumColumns; column++ {
		for row := 0; row <= (NumRows - NumConsecutive); row++ {
			if board.checkSouthwestDiagonalWin(row, column) {
				return true
			}
		}
	}
	return false
}

func (board *Board) ToStringBoard() [NumRows][NumColumns]string {
	var stringBoard [NumRows][NumColumns]string
	for row := 0; row < NumRows; row++ {
		for column := 0; column < NumColumns; column++ {
			if board[row][column] == Empty {
				stringBoard[row][column] = ""
			} else if board[row][column] == Red {
				stringBoard[row][column] = "R"
			} else if board[row][column] == Black {
				stringBoard[row][column] = "B"
			} else {
				panic(fmt.Sprintf("invalid value", board[row][column], "for a board"))
			}
		}
	}
	return stringBoard
}
func (board *Board) IsFull() bool {
	// will check the top row, which is always the last to fill up.
	for column := 0; column < NumColumns; column++ {
		if board[0][column] == Empty {
			return false
		}
	}
	return true
}

// Returns error if the move is invalid
func (bp *Board) ApplyMove(move int, playerId int) error {
	if move >= NumColumns || move < 0 {
		return errors.New(fmt.Sprintf("Move %d is invalid", move))
	}
	for i := NumRows - 1; i >= 0; i-- {
		if bp[i][move] == 0 {
			bp[i][move] = playerId
			return nil
		}
	}
	return errors.New(fmt.Sprintf("No room in column %d for a move", move))
}

func NewBoard() *Board {
	// Board is initialized to be filled with zeros.
	var board Board
	return &board
}
