package tictactoe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestHandlers(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			g := NewGame(random, greedy)
			g.Start()
			t.Logf("%+v", g.Board)
			wg.Done()
		}()
	}
	wg.Wait()
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

func TestWebhookHandle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload HookPayload
		var move HookMove
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Println(string(body))
		err = json.Unmarshal(body, &payload)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		space, err := random(payload.Piece, payload.Board)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		move.Space = space
		blob, err := json.MarshalIndent(move, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, string(blob))
	}))
	defer ts.Close()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			g := NewGame(NewWebhookHandlerFunc(ts.URL), NewWebhookHandlerFunc(ts.URL))
			g.Start()
			t.Logf("%+v", g.Board)
			wg.Done()
		}()
	}
	wg.Wait()
}
