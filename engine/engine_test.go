package engine

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestMatch struct {
	Current Player
}

func (match *TestMatch) CurrentPlayer() Player {
	return match.Current
}

func (match *TestMatch) Play(p Player, body []byte) (bool, error) {
	return string(body) == "hello\n", nil
}

func (match *TestMatch) Winner() (Player, error) {
	return Player{}, nil
}

func (match *TestMatch) Stalemate() bool {
	return false
}


func TestUnreachableServer(t *testing.T) {
	t.Parallel()

	match := TestMatch{}
	err := PlayMatch(&match)

	if err == nil {
		t.Fatalf("Server should have errored")
	}
}


func TestWinningMove(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method to be POST, instead was %s", r.Method)
		}
		fmt.Fprintln(w, "hello")
	}))

	match := TestMatch{Current: Player{MatchUrl: ts.URL}}

	err := PlayMatch(&match)


	if err != nil {
		t.Logf("err:", err)
		t.Fatalf("Server shouldn't have errored")
	}
}
