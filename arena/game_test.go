package main

import (
	"testing"
)

func TestGameOver(t *testing.T) {
	winVertical := [7][7]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
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

	winOtherVertical := [7][7]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
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
	winHorizontal := [7][7]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
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

	winDiagonal := [7][7]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
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

	winSouthwestDiagonal := [7][7]int{
		[7]int{0, 0, 0, 0, 0, 0, 0},
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

	unfinishedGame := [7][7]int{
		[7]int{2, 0, 0, 1, 0, 0, 0},
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
