package tictactoe

import (
	"fmt"
	"log"
)

type Piece int
type Board [9]Piece

const (
	X Piece = 1
	O Piece = 2
)

func (p Piece) String() string {
	switch p {
	case X:
		return "X"
	case O:
		return "O"
	}
	return " "
}

func (b Board) String() string {
	repr := `
    %s|%s|%s
    -----
    %s|%s|%s
    -----
    %s|%s|%s
    `
	return fmt.Sprintf(repr, b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8])
}

// A tic-tac-toe board is represented as nine-length array of integers
type Game struct {
	Board   Board
	Players []Player
}

type Player struct {
	Piece   Piece
	Actions chan int
	Updates chan Board
}

func (p Player) Play() {
	for {
		board, ok := <-p.Updates

		if !ok {
			close(p.Actions)
			return
		}

		for i, slot := range board {
			if slot == 0 {
				p.Actions <- i
				break
			}
		}
	}
}

func (g *Game) Add(piece Piece, i int) error {
	if i < 0 || i > 8 {
		return fmt.Errorf("square %d is out of bounds", i)
	}
	if g.Board[i] != 0 {
		return fmt.Errorf("square %s is occupied by %d", i, g.Board[i])
	}
	g.Board[i] = piece
	return nil
}

func (g *Game) Start() {
	for _, p := range g.Players {
		go p.Play()
	}

	takeTurn := func(p Player) bool {
		p.Updates <- g.Board

		square, ok := <-p.Actions

		if !ok {
			return true
		}

		fmt.Printf("Player %d picked spot %d\n", p.Piece, square)

		err := g.Add(p.Piece, square)

		if err != nil {
			//Fixme
			log.Fatal(err)
		}

		if g.IsOver() {
			fmt.Println("Game over")
			for _, p := range g.Players {
				close(p.Updates)
			}
			return true
		}

		return false
	}

	loop := true

	for loop {
		for _, p := range g.Players {
			if !loop {
				continue
			}
			if takeTurn(p) {
				loop = false
				break
			}
		}
	}

	for _, p := range g.Players {
		<-p.Actions
	}
}

func (g Game) IsOver() bool {
	over := func(i, j, k int) bool {
		return g.Board[i] == g.Board[j] &&
			g.Board[j] == g.Board[k] &&
			g.Board[i] != 0
	}

	wins := []bool{
		// Rows
		over(0, 1, 2),
		over(3, 4, 5),
		over(6, 7, 8),
		// Columns
		over(0, 3, 6),
		over(1, 4, 7),
		over(2, 5, 8),
		// Diag
		over(0, 4, 8),
		over(2, 4, 6),
	}

	for _, win := range wins {
		if win {
			return true
		}
	}

	for _, piece := range g.Board {
		if piece == 0 {
			return false
		}
	}

	return false
}

func NewPlayer(piece Piece) Player {
	return Player{Piece: piece, Actions: make(chan int), Updates: make(chan Board)}
}

func NewGame() Game {
	return Game{Board: Board{}, Players: []Player{NewPlayer(O), NewPlayer(X)}}
}
