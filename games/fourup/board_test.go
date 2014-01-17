package fourup

import (
	"testing"
)

func TestBoardFull(t *testing.T) {
	t.Parallel()

	fullBoard := Board{
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
	}

	if !fullBoard.IsFull() {
		t.Errorf("Full board should be marked full")
	}

	boardWithRoom := Board{
		[7]int8{2, 0, 2, 2, 2, 2, 2},
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
	}

	if boardWithRoom.IsFull() {
		t.Errorf("Board with room be marked not full")
	}
}

func TestGameOver(t *testing.T) {
	t.Parallel()

	winThirdVertical := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{1, 0, 0, 0, 0, 0, 1},
		[7]int8{1, 0, 0, 0, 0, 0, 2},
		[7]int8{1, 0, 0, 0, 0, 0, 2},
	}
	if over, _ := winThirdVertical.GameOver(); !over {
		t.Errorf("Game should be over if 4 vertical tiles " +
			"starting in top row, form a connect four")
	}

	winVertical := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
	}
	if over, _ := winVertical.GameOver(); !over {
		t.Errorf("Game should be over if 4 vertical tiles are in a row")
	}

	winOtherVertical := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
	}
	if over, _ := winOtherVertical.GameOver(); !over {
		t.Errorf("Game should be over if 4 other vertical tiles are in a row")
	}

	winHorizontal := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 1, 1, 1, 1, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
	}
	if over, _ := winHorizontal.GameOver(); !over {
		t.Errorf("Game should be over if 4 horizontal tiles are in a row")
	}

	winDiagonal := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 1, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 1, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 1, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 1, 0},
	}
	if over, _ := winDiagonal.GameOver(); !over {
		t.Errorf("Game should be over if 4 southeast diagonal tiles are in a row")
	}

	winSouthwestDiagonal := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 1, 0, 0},
		[7]int8{0, 0, 0, 1, 0, 0, 0},
		[7]int8{0, 0, 1, 0, 1, 0, 0},
		[7]int8{0, 1, 0, 0, 0, 1, 0},
	}
	if over, _ := winSouthwestDiagonal.GameOver(); !over {
		t.Errorf("Game should be over if 4 southwest diagonal tiles are in a row")
	}
	if _, winner := winSouthwestDiagonal.GameOver(); winner != 1 {
		t.Errorf("winner should be 1")
	}

	unfinishedGame := Board{
		[7]int8{0, 2, 0, 2, 0, 0, 0},
		[7]int8{0, 0, 2, 2, 1, 1, 1},
		[7]int8{0, 0, 1, 1, 2, 2, 2},
		[7]int8{0, 0, 2, 1, 2, 1, 0},
		[7]int8{0, 0, 0, 2, 1, 1, 0},
		[7]int8{0, 0, 1, 1, 2, 2, 2},
	}
	if over, _ := unfinishedGame.GameOver(); over {
		t.Errorf("Game was marked over, but wasn't over")
	}
}

func TestApplyMoveToBoard(t *testing.T) {
	t.Parallel()
	emptyBoard := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
	}

	oneMoveBoard := Board{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
	}

	_ = emptyBoard.ApplyMove(6, 1)

	if emptyBoard != oneMoveBoard {
		t.Errorf("New board does not equal board with expected move")
	}

	columnFullBoard := Board{
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
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
	t.Parallel()
	fullBoard := Board{
		[7]int8{1, 2, 1, 2, 0, 0, 0},
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
	}
	expectedBoard := StringBoard{
		[7]string{"R", "B", "R", "B", "", "", ""},
		[7]string{"B", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
	}
	stringBoard := fullBoard.getStringBoard()
	if stringBoard != expectedBoard {
		t.Fatalf("Output board", stringBoard, "didn't match expected board")
	}
}
