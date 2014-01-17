package fourup

import (
	"testing"
)

func TestBoardFull(t *testing.T) {
	t.Parallel()
	fullBoard := [NumRows][NumColumns]int8{
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
	}
	if !IsBoardFull(fullBoard) {
		t.Errorf("Full board should be marked full")
	}

	boardWithRoom := [NumRows][NumColumns]int8{
		[7]int8{2, 0, 2, 2, 2, 2, 2},
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
	}
	if IsBoardFull(boardWithRoom) {
		t.Errorf("Board with room be marked not full")
	}
}

func TestGameOver(t *testing.T) {
	t.Parallel()

	winThirdVertical := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{1, 0, 0, 0, 0, 0, 1},
		[7]int8{1, 0, 0, 0, 0, 0, 2},
		[7]int8{1, 0, 0, 0, 0, 0, 2},
	}
	if !GameOver(winThirdVertical) {
		t.Errorf("Game should be over if 4 vertical tiles " +
			"starting in top row, form a connect four")
	}

	winVertical := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
	}
	if !GameOver(winVertical) {
		t.Errorf("Game should be over if 4 vertical tiles are in a row")
	}

	winOtherVertical := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
	}
	if !GameOver(winOtherVertical) {
		t.Errorf("Game should be over if 4 other vertical tiles are in a row")
	}

	winHorizontal := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 1, 1, 1, 1, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
	}
	if !GameOver(winHorizontal) {
		t.Errorf("Game should be over if 4 horizontal tiles are in a row")
	}

	winDiagonal := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 1, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 1, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 1, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 1, 0},
	}
	if !GameOver(winDiagonal) {
		t.Errorf("Game should be over if 4 southeast diagonal tiles are in a row")
	}

	winSouthwestDiagonal := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 1, 0, 0},
		[7]int8{0, 0, 0, 1, 0, 0, 0},
		[7]int8{0, 0, 1, 0, 1, 0, 0},
		[7]int8{0, 1, 0, 0, 0, 1, 0},
	}
	if !GameOver(winSouthwestDiagonal) {
		t.Errorf("Game should be over if 4 southwest diagonal tiles are in a row")
	}

	unfinishedGame := [NumRows][NumColumns]int8{
		[7]int8{0, 2, 0, 2, 0, 0, 0},
		[7]int8{0, 0, 2, 2, 1, 1, 1},
		[7]int8{0, 0, 1, 1, 2, 2, 2},
		[7]int8{0, 0, 2, 1, 2, 1, 0},
		[7]int8{0, 0, 0, 2, 1, 1, 0},
		[7]int8{0, 0, 1, 1, 2, 2, 2},
	}
	if GameOver(unfinishedGame) {
		t.Errorf("Game was marked over, but wasn't over")
	}
}

func TestApplyMoveToBoard(t *testing.T) {
	t.Parallel()
	emptyBoard := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
	}

	oneMoveBoard := [NumRows][NumColumns]int8{
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 0},
		[7]int8{0, 0, 0, 0, 0, 0, 1},
	}

	newBoard, _ := ApplyMoveToBoard(6, 1, &emptyBoard)
	if *newBoard != oneMoveBoard {
		t.Errorf("New board does not equal board with expected move")
	}

	columnFullBoard := [NumRows][NumColumns]int8{
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
		[7]int8{1, 0, 0, 0, 0, 0, 0},
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

func TestStringBoard(t *testing.T) {
	t.Parallel()
	fullBoard := [NumRows][NumColumns]int8{
		[7]int8{1, 2, 1, 2, 0, 0, 0},
		[7]int8{2, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
		[7]int8{1, 2, 2, 2, 2, 2, 2},
	}
	expectedBoard := [NumRows][NumColumns]string{
		[7]string{"R", "B", "R", "B", "", "", ""},
		[7]string{"B", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
		[7]string{"R", "B", "B", "B", "B", "B", "B"},
	}
	stringBoard := GetStringBoard(&fullBoard)
	if stringBoard != expectedBoard {
		t.Fatalf("Output board", stringBoard, "didn't match expected board")
	}
}

func TestGetIntBoard(t *testing.T) {
	t.Parallel()
	dbBoard := []byte(`[["B","","","R","","","B"],["R","","","B","","","R"],
	["R","","B","R","B","","R"],["R","R","B","R","B","R","B"],
	["B","R","B","B","R","B","R"],["B","R","B","B","R","B","R"]]`)
	expectedBoard := &[NumRows][NumColumns]int8{
		[7]int8{Black, 0, 0, Red, 0, 0, Black},
		[7]int8{Red, 0, 0, Black, 0, 0, Red},
		[7]int8{Red, 0, Black, Red, Black, 0, Red},
		[7]int8{Red, Red, Black, Red, Black, Red, Black},
		[7]int8{Black, Red, Black, Black, Red, Black, Red},
		[7]int8{Black, Red, Black, Black, Red, Black, Red},
	}
	board, err := GetIntBoard(dbBoard)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if *board != *expectedBoard {
		t.Fatalf("Output board %d, didn't match expected board", board)
	}
}
