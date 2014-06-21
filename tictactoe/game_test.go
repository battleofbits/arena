package tictactoe

import (
	"testing"
)

func TestHandlers(t *testing.T) {
	for i := 0; i < 2; i++ {
		g := NewGame(random, greedy)
		g.Start()
		t.Logf("%+v", g.Board)
	}
}

func TestWinConditions(t *testing.T) {
	wins := []Board{
		Board{
			X, 0, 0,
			X, 0, 0,
			X, 0, 0,
		},
		Board{
			X, X, X,
			0, 0, 0,
			0, 0, 0,
		},
		Board{
			X, O, 0,
			0, X, 0,
			0, O, X,
		},
	}
	for _, win := range wins {
		g := Game{Board: win}
		if !g.IsOver() {
			t.Errorf("Game should be over %s", win)
		}
	}
}
