package tictactoe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type Piece int
type Board [9]Piece
type HandlerFunc func(Piece, Board) (int, error)

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
type Player struct {
	Piece   Piece
	Actions chan int
	Updates chan Board
	Logic   HandlerFunc
}

func (p Player) Play() {
	for {
		board, ok := <-p.Updates
		if !ok {
			close(p.Actions)
			return
		}
		i, err := p.Logic(p.Piece, board)
		if err != nil {
			log.Println(err)
			close(p.Actions)
			return
		}
		p.Actions <- i
	}
}

// Game Handlers
func random(p Piece, board Board) (int, error) {
	for i := 0; i < 20; i++ {
		if i := rand.Intn(9); board[i] == 0 {
			return i, nil
		}
	}
	return 0, fmt.Errorf("tried 20 random spaces")
}

func greedy(p Piece, board Board) (int, error) {
	for i, slot := range board {
		if slot == 0 {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no open spaces")
}

type HookPayload struct {
	Piece Piece `json:"piece"`
	Board Board `json:"board"`
}

type HookMove struct {
	Space int `json:"space"`
}

func NewWebhookHandlerFunc(url string) HandlerFunc {
	return func(p Piece, board Board) (int, error) {
		payload := HookPayload{Piece: p, Board: board}
		blob, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return 0, err
		}
		resp, err := http.Post(url, "application/json", bytes.NewReader(blob))
		if err != nil {
			return 0, err
		}
		if resp.StatusCode >= http.StatusBadRequest {
			return 0, fmt.Errorf("hook %s returned status %d", url, resp.StatusCode)
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		var move HookMove
		err = json.Unmarshal(data, &move)
		if err != nil {
			return 0, err
		}
		return move.Space, nil
	}
}

type Game struct {
	Board   Board
	Players []Player
}

func (g *Game) Add(piece Piece, i int) error {
	if i < 0 || i > 8 {
		return fmt.Errorf("space %d is out of bounds", i)
	}
	if g.Board[i] != 0 {
		return fmt.Errorf("space %s is occupied by %d", i, g.Board[i])
	}
	g.Board[i] = piece
	return nil
}

// Zero tolerance for invalid moves
func (g *Game) TakeTurn(p Player) error {
	p.Updates <- g.Board
	square, ok := <-p.Actions
	if !ok {
		return fmt.Errorf("action channel closed")
	}
	return g.Add(p.Piece, square)
}

func (g *Game) Start() error {
	for _, p := range g.Players {
		go p.Play()
	}
	func() {
		for {
			for _, p := range g.Players {
				err := g.TakeTurn(p)
				if err != nil || g.IsOver() {
					return
				}
			}
		}
	}()
	for _, p := range g.Players {
		close(p.Updates)
	}
	return nil
}

func (g Game) IsOver() bool {
	match := func(i, j, k int) bool {
		return g.Board[i] == g.Board[j] &&
			g.Board[j] == g.Board[k] &&
			g.Board[i] != 0
	}
	wins := []bool{
		match(0, 1, 2), // horizontal rows
		match(3, 4, 5),
		match(6, 7, 8),
		match(0, 3, 6), // vertical rows
		match(1, 4, 7),
		match(2, 5, 8),
		match(0, 4, 8), // diagonal rows
		match(2, 4, 6),
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

func newPlayer(piece Piece, f HandlerFunc) Player {
	return Player{
		Piece:   piece,
		Actions: make(chan int),
		Updates: make(chan Board),
		Logic:   f,
	}
}

func NewGame(one, two HandlerFunc) Game {
	return Game{Board: Board{},
		Players: []Player{newPlayer(X, one), newPlayer(O, two)}}
}
