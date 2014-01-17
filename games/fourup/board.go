package fourup

import (
	"encoding/json"
	"errors"
	"fmt"
)

const NumRows = 6
const NumColumns = 7
const NumConsecutive = 4

// Checks fourup board state
// Author: Kevin Burke <kev@inburke.com>

// row varies, column does not.
func checkVerticalWin(column int8, board [NumRows][NumColumns]int8) (bool, int8) {
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

func checkHorizontalWin(row int8, board [NumRows][NumColumns]int8) (bool, int8) {
	checkColumnInRow := func(row int8, column int8,
		board [NumRows][NumColumns]int8) (bool, int8) {
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
func checkSoutheastDiagonalWin(row int8, column int8,
	board [NumRows][NumColumns]int8) (bool, int8) {

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

func checkSouthwestDiagonalWin(row int8, column int8,
	board [NumRows][NumColumns]int8) (bool, int8) {

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
func gameOver(board [NumRows][NumColumns]int8) (bool, int8) {
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

func GetStringBoard(board *[NumRows][NumColumns]int8) [NumRows][NumColumns]string {
	var stringBoard [NumRows][NumColumns]string
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

// Convert the database representation of a board into a regular board
func GetIntBoard(dbBoard []byte) (*[NumRows][NumColumns]int8, error) {
	var stringBoard [NumRows][NumColumns]string
	err := json.Unmarshal(dbBoard, &stringBoard)
	if err != nil {
		return nil, err
	}

	var board [NumRows][NumColumns]int8
	for row := int8(0); row < NumRows; row++ {
		for column := int8(0); column < NumColumns; column++ {
			if stringBoard[row][column] == "" {
				board[row][column] = Empty
			} else if stringBoard[row][column] == "R" {
				board[row][column] = Red
			} else if stringBoard[row][column] == "B" {
				board[row][column] = Black
			} else {
				panic(fmt.Sprint("invalid value ", stringBoard[row][column],
					" for a board"))
			}
		}
	}
	return &board, nil
}

func IsBoardFull(board [NumRows][NumColumns]int8) bool {
	// will check the top row, which is always the last to fill up.
	for column := 0; column < NumColumns; column++ {
		if board[0][column] == Empty {
			return false
		}
	}
	return true
}

// Returns error if the move is invalid
func applyMoveToBoard(move int8, color int8, bp *[NumRows][NumColumns]int8) (
	*[NumRows][NumColumns]int8, error) {
	if move >= NumColumns || move < 0 {
		return bp, errors.New(fmt.Sprintf("Move %d is invalid", move))
	}
	for i := NumRows - 1; i >= 0; i-- {
		if bp[i][move] == 0 {
			bp[i][move] = color
			return bp, nil
		}
	}
	return bp, errors.New(fmt.Sprintf("No room in column %d for a move", move))
}

func initializeBoard() *[NumRows][NumColumns]int8 {
	// Board is initialized to be filled with zeros.
	var board [NumRows][NumColumns]int8
	return &board
}
