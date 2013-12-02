package fourup

import (
	"testing"
)

func TestBoardFull(t *testing.T) {
	fullBoard := Board{
		[7]int{2, 2, 2, 2, 2, 2, 2},
		[7]int{2, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
	}
	if !fullBoard.IsFull() {
		t.Errorf("Full board should be marked full")
	}

	boardWithRoom := Board{
		[7]int{2, 0, 2, 2, 2, 2, 2},
		[7]int{2, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
	}
	if boardWithRoom.IsFull() {
		t.Errorf("Board with room be marked not full")
	}
}

func TestGameOver(t *testing.T) {

	winThirdVertical := Board{
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{1, 0, 0, 0, 0, 0, 1},
		[7]int{1, 0, 0, 0, 0, 0, 2},
		[7]int{1, 0, 0, 0, 0, 0, 2},
	}
	if !winThirdVertical.GameOver() {
		t.Errorf("Game should be over if 4 vertical tiles starting in top row, form a connect four")
	}

	winVertical := Board{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
	}
	if !winVertical.GameOver() {
		t.Errorf("Game should be over if 4 vertical tiles are in a row")
	}

	winOtherVertical := Board{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
		[7]int{0, 0, 0, 0, 0, 0, 1},
	}
	if !winOtherVertical.GameOver() {
		t.Errorf("Game should be over if 4 other vertical tiles are in a row")
	}

	winHorizontal := Board{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 1, 1, 1, 1, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
	}
	if !winHorizontal.GameOver() {
		t.Errorf("Game should be over if 4 horizontal tiles are in a row")
	}

	winDiagonal := Board{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 1, 0, 0, 0, 0},
		[7]int{0, 0, 0, 1, 0, 0, 0},
		[7]int{0, 0, 0, 0, 1, 0, 0},
		[7]int{0, 0, 0, 0, 0, 1, 0},
	}
	if !winDiagonal.GameOver() {
		t.Errorf("Game should be over if 4 southeast diagonal tiles are in a row")
	}

	winSouthwestDiagonal := Board{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 1, 0, 0},
		[7]int{0, 0, 0, 1, 0, 0, 0},
		[7]int{0, 0, 1, 0, 1, 0, 0},
		[7]int{0, 1, 0, 0, 0, 1, 0},
	}
	if !winSouthwestDiagonal.GameOver() {
		t.Errorf("Game should be over if 4 southwest diagonal tiles are in a row")
	}

	unfinishedBoard := Board{
		[7]int{0, 2, 0, 2, 0, 0, 0},
		[7]int{0, 0, 2, 2, 1, 1, 1},
		[7]int{0, 0, 1, 1, 2, 2, 2},
		[7]int{0, 0, 2, 1, 2, 1, 0},
		[7]int{0, 0, 0, 2, 1, 1, 0},
		[7]int{0, 0, 1, 1, 2, 2, 2},
	}
	if unfinishedBoard.GameOver() {
		t.Errorf("Game was marked over, but wasn't over")
	}
}

func TestApplyMoveToBoard(t *testing.T) {
	emptyBoard := Board{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
	}

	oneMoveBoard := Board{
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 0},
		[7]int{0, 0, 0, 0, 0, 0, 1},
	}

	emptyBoard.ApplyMove(6, 1)

	if emptyBoard != oneMoveBoard {
		t.Errorf("New board does not equal board with expected move")
	}

	columnFullBoard := Board{
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
		[7]int{1, 0, 0, 0, 0, 0, 0},
	}

	err := columnFullBoard.ApplyMove(0, 1)

	if err.Error() != "No room in column 0 for a move" {
		t.Errorf("Should have rejected move in column 0, did not, error was %s", err.Error())
	}

	err = columnFullBoard.ApplyMove(-22, 1)
	if err == nil || err.Error() != "Move -22 is invalid" {
		t.Errorf("Should have rejected negative move, did not, error was %s", err.Error())
	}

	err = columnFullBoard.ApplyMove(7, 1)
	if err == nil || err.Error() != "Move 7 is invalid" {
		t.Errorf("Should have rejected positive move, did not, error was %s", err.Error())
	}
}

func TestStringBoard(t *testing.T) {
	fullBoard := Board{
		[7]int{1, 2, 1, 2, 0, 0, 0},
		[7]int{2, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
		[7]int{1, 2, 2, 2, 2, 2, 2},
	}
	expectedBoard := [NumRows][NumColumns]string{
		[7]string{"R", "B", "R", "B", "", "", ""},
		[7]string{"B", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
	}
	stringBoard := fullBoard.ToStringBoard()
	if stringBoard != expectedBoard {
		t.Errorf("Output board", stringBoard, "didn't match expected board")
	}
}
