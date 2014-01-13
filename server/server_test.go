package main

import (
	"encoding/json"
	"fmt"
	"github.com/battleofbits/arena/arena"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	//"reflect"
	"strings"
	"testing"
	"time"
)

func TestMoves(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/games/four-up/matches/{match}/moves", MovesHandler)
	req, _ := http.NewRequest("GET", "http://localhost/games/four-up/matches/3/moves", nil)

	player := "kevin"
	column := 3
	now := time.Now()
	move := &Move{
		Player: player,
		Column: column,
		Played: now,
	}

	moveGetter = func(moveId int) []*Move {
		return []*Move{move}
	}

	defer reassignMoveGetter(getMoves)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	var moves Moves
	err := json.Unmarshal(resp.Body.Bytes(), &moves)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(moves.Moves) != 1 {
		t.Errorf("expected only one result back but got %d", len(moves.Moves))
	}
	endMove := moves.Moves[0]
	if endMove.Player != player {
		t.Errorf("expected player to be %s, got %s", player, endMove.Player)
	}
	if endMove.Column != column {
		t.Errorf("expected column to be %d, got %d", column, endMove.Column)
	}
	// For some reason you can't compare the timestamps with Equal, it doesn't
	// work. No idea why.
	if endMove.Played.Unix() != now.Unix() {
		t.Errorf("expected timestamp to be %s, got %s", now, endMove.Played)
	}
}

func reassignMoveGetter(to func(int) []*Move) {
	moveGetter = to
}

func TestEmptyMatches(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/games/four-up/matches", MatchesHandler)
	req, _ := http.NewRequest("GET", "http://localhost/games/four-up/matches", nil)

	matchesGetter = func() ([]*arena.FourUpMatch, error) {
		return []*arena.FourUpMatch{}, nil
	}

	defer reassignMatchesGetter(arena.GetMatches)

	resp := httptest.NewRecorder()
	var response Response
	r.ServeHTTP(resp, req)
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(err.Error())
	}
	matches := response["matches"].([]interface{})
	if len(matches) != 0 {
		t.Fatalf("match length should have been 0, was %d", len(matches))
	}
}

type MatchResponses struct {
	Matches []arena.MatchResponse `json:"matches"`
}

func TestOneMatch(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/games/four-up/matches", MatchesHandler)
	req, _ := http.NewRequest("GET", "http://localhost/games/four-up/matches", nil)

	match := &arena.FourUpMatch{
		Id:     3,
		Winner: nil,
	}
	matchesGetter = func() ([]*arena.FourUpMatch, error) {
		return []*arena.FourUpMatch{match}, nil
	}

	defer reassignMatchesGetter(arena.GetMatches)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var response MatchResponses
	bits := resp.Body.Bytes()
	fmt.Println(string(bits))
	err := json.Unmarshal(bits, &response)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(response.Matches) != 1 {
		t.Fatalf("match length should have been 0, was %d", len(response.Matches))
	}
	mr := response.Matches[0]
	fmt.Println(mr.Winner)
	if mr.Id != 3 {
		t.Errorf("id should have been 3, was %d", mr.Id)
	}
}

func reassignMatchGetter(to func(int) (*arena.FourUpMatch, error)) {
	matchGetter = to
}

func reassignMatchesGetter(to func() ([]*arena.FourUpMatch, error)) {
	matchesGetter = to
}

func TestSendInviteOK(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method to be POST, instead was %s", r.Method)
		}
		decoder := json.NewDecoder(r.Body)
		var is *InviteRequest
		err := decoder.Decode(&is)
		if err != nil {
			t.Errorf(err.Error())
		}
		fmt.Println(is.Game)
		if is.Game != "fourup" {
			t.Errorf("expected 'game' to be fourup, instead was %s", is.Game)
		}
		if is.RequestingPlayer != "kevinburke" {
			t.Errorf("expected requesting player to be kevinburke, "+
				"instead was %s", is.RequestingPlayer)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, client")
	}))

	err := SendInvite(ts.URL, "fourup", "kevinburke")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSendInviteError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Hello, client")
	}))
	err := SendInvite(ts.URL, "fourup", "kevinburke")
	if err != nil {
		expected := "Received error status 400 Bad Request"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("got weird error %s, expected '%s'", err.Error(), expected)
		}
	}
}
