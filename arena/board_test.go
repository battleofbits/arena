package main

import (
	"testing"
)

func testBoardFull(t *testing.T) {
	fullBoard := [6][7]int{
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

	boardWithRoom := [6][7]int{
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
	winVertical := [6][7]int{
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

	winOtherVertical := [6][7]int{
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
	winHorizontal := [6][7]int{
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

	winDiagonal := [6][7]int{
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

	winSouthwestDiagonal := [6][7]int{
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

	unfinishedGame := [6][7]int{
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
	emptyBoard := [6][7]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
	}

	oneMoveBoard := [6][7]int{
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

	columnFullBoard := [6][7]int{
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
