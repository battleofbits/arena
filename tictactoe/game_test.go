package tictactoe

import (
	"testing"
)

func TestGame(t *testing.T) {
	g := NewGame()
	g.Start()
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
