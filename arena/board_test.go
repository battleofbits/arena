package main

import (
	"testing"
)

func TestBoardFull(t *testing.T) {
	fullBoard := [NumRows][NumColumns]int{
		[7]int{2, 2, 2, 2, 2, 2, 2},
		[7]int{2, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
	}
	if !IsBoardFull(fullBoard) {
		t.Errorf("Full board should be marked full")
	}

	boardWithRoom := [NumRows][NumColumns]int{
		[7]int{2, 0, 2, 2, 2, 2, 2},
		[7]int{2, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
	}
	if IsBoardFull(boardWithRoom) {
		t.Errorf("Board with room be marked not full")
	}
}

func TestGameOver(t *testing.T) {

	winThirdVertical := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{1, 0, 0, 0, 0, 0, 1},
		[7]int{1, 0, 0, 0, 0, 0, 2},
		[7]int{1, 0, 0, 0, 0, 0, 2},
	}
	if !GameOver(winThirdVertical) {
		t.Errorf("Game should be over if 4 vertical tiles starting in top row, form a connect four")
	}

	winVertical := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
	}
	if !GameOver(winVertical) {
		t.Errorf("Game should be over if 4 vertical tiles are in a row")
	}

	winOtherVertical := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
	}
	if !GameOver(winOtherVertical) {
		t.Errorf("Game should be over if 4 other vertical tiles are in a row")
	}

	winHorizontal := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 1, 1, 1, 1, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
	}
	if !GameOver(winHorizontal) {
		t.Errorf("Game should be over if 4 horizontal tiles are in a row")
	}

	winDiagonal := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 1, 0, 0, 0, 0},
		[7]int{0, 0, 0, 1, 0, 0, 0},
		[7]int{0, 0, 0, 0, 1, 0, 0},
		[7]int{0, 0, 0, 0, 0, 1, 0},
	}
	if !GameOver(winDiagonal) {
		t.Errorf("Game should be over if 4 southeast diagonal tiles are in a row")
	}

	winSouthwestDiagonal := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 1, 0, 0},
		[7]int{0, 0, 0, 1, 0, 0, 0},
		[7]int{0, 0, 1, 0, 1, 0, 0},
		[7]int{0, 1, 0, 0, 0, 1, 0},
	}
	if !GameOver(winSouthwestDiagonal) {
		t.Errorf("Game should be over if 4 southwest diagonal tiles are in a row")
	}

	unfinishedGame := [NumRows][NumColumns]int{
		[7]int{0, 2, 0, 2, 0, 0, 0},
		[7]int{0, 0, 2, 2, 1, 1, 1},
		[7]int{0, 0, 1, 1, 2, 2, 2},
		[7]int{0, 0, 2, 1, 2, 1, 0},
		[7]int{0, 0, 0, 2, 1, 1, 0},
		[7]int{0, 0, 1, 1, 2, 2, 2},
	}
	if GameOver(unfinishedGame) {
		t.Errorf("Game was marked over, but wasn't over")
	}
}

func TestApplyMoveToBoard(t *testing.T) {
	emptyBoard := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
	}

	oneMoveBoard := [NumRows][NumColumns]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 1},
	}

	newBoard, _ := ApplyMoveToBoard(6, 1, &emptyBoard)
	if *newBoard != oneMoveBoard {
		t.Errorf("New board does not equal board with expected move")
	}

	columnFullBoard := [NumRows][NumColumns]int{
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
	}

	_, err := ApplyMoveToBoard(0, 1, &columnFullBoard)
	if err.Error() != "No room in column 0 for a move" {
		t.Errorf("Should have rejected move in column 0, did not, error was %s", err.Error())
	}

	_, err = ApplyMoveToBoard(-22, 1, &columnFullBoard)
	if err == nil || err.Error() != "Move -22 is invalid" {
		t.Errorf("Should have rejected negative move, did not, error was %s", err.Error())
	}

	_, err = ApplyMoveToBoard(7, 1, &columnFullBoard)
	if err == nil || err.Error() != "Move 7 is invalid" {
		t.Errorf("Should have rejected positive move, did not, error was %s", err.Error())
	}
}
