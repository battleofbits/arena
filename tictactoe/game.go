package main

import (
	"fmt"
)

type Player struct {
	Piece   int
	Actions chan int
	Updates chan [9]int // Come up with a better name
}

type Game struct {
	Board     [9]int
	PlayerOne Player
	PlayerTwo Player
	PlayerThree Player
}

func (g Game) Start() {
	players := []Player{g.PlayerOne, g.PlayerTwo, g.PlayerThree}

	play := func(player Player) {
		for {
			board, ok := <-player.Updates

			if !ok {
				close(player.Actions)
				return
			}

			for i, slot := range board {
				if slot == 0 {
					player.Actions <- i
					fmt.Printf("Player %d picked spot %d\n", player.Piece, i)
					break
				}
			}
		}
	}

	for _, p := range players {
		go play(p)
	}

	takeTurn := func(p Player) bool {
		p.Updates <- g.Board

		move, ok := <-p.Actions

		if !ok {
			return true
		}

		g.Board[move] = p.Piece

		if g.IsOver() {
			fmt.Println("Game over")
			for _, p := range players {
				close(p.Updates)
			}
			return true
		}

		return false
	}

	loop := true

	for loop {
		for _, p := range players {
			if !loop {
				continue
			}
			if takeTurn(p) {
				loop = false
				break
			}
		}
	}

	for _, p := range players {
		<-p.Actions
	}
}

func (g Game) IsOver() bool {
	for _, piece := range g.Board {
		if piece == 0 {
			return false
		}
	}
	return true
}

func NewPlayer(piece int) Player {
	return Player{Piece: piece, Actions: make(chan int), Updates: make(chan [9]int)}
}

func NewGame() Game {
        return Game{Board: [9]int{}, PlayerOne: NewPlayer(1), PlayerTwo: NewPlayer(2), PlayerThree: NewPlayer(4)}
}

func main() {
	g := NewGame()
	g.Start()
}
