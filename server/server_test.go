package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/battleofbits/arena/arena"
	"github.com/battleofbits/arena/engine"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	//"reflect"
	//"strings"
	"testing"
	"time"
)

// Retrieving a list of moves should return the expected response.
// This mocks out the database call to fetch the list of moves.
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

// Sending an invitation without specifying a game should return a 400 Bad
// Request.
func TestInviteNoGame(t *testing.T) {
	t.Parallel()

	r := mux.NewRouter()
	buf := bytes.NewBufferString("{}")
	r.HandleFunc("/players/{player}/invitations", InvitationsHandler)
	req, _ := http.NewRequest("POST", "http://localhost/players/kevinburke/invitations", buf)
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != 400 {
		t.Errorf("Expected status 400 but got %d", resp.Code)
	}
	var err Error
	decodingErr := json.Unmarshal(resp.Body.Bytes(), &err)
	if decodingErr != nil {
		fmt.Println(string(resp.Body.Bytes()))
		t.Fatalf(decodingErr.Error())
	}
	if err.Message != "No game specified" {
		t.Errorf("Expected error message to read 'No game specified', was '%s'", err.Message)
	}
	if err.Type != "invalid-game" {
		t.Errorf("Expected error type to be 'invalid-game', was '%s'", err.Type)
	}
}

// Sending an invitation to an unknown player should return a 400 Bad Request.
func TestInviteUnknownPlayer(t *testing.T) {

	// Reassign the player getter function to return no database rows.
	getDatastore = func() engine.Datastore {
		return engine.NotFoundDatastore{}
	}
	defer reassignDatastoreGetter()

	r := mux.NewRouter()
	buf := bytes.NewBufferString("{\"Game\": \"fourup\"}")
	r.HandleFunc("/players/{player}/invitations", InvitationsHandler)
	req, _ := http.NewRequest("POST", "http://localhost/players/foobar/invitations", buf)
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != 404 {
		t.Errorf("Expected status 404 but got %d", resp.Code)
	}
	var err Error
	decodingErr := json.Unmarshal(resp.Body.Bytes(), &err)
	if decodingErr != nil {
		fmt.Println(string(resp.Body.Bytes()))
		t.Fatalf(decodingErr.Error())
	}
	playerErrMsg := "No players with name foobar"
	playerErrType := "invalid-player"
	if err.Message != playerErrMsg {
		t.Errorf("Expected error message to read '%s', was '%s'", playerErrMsg, err.Message)
	}
	if err.Type != playerErrType {
		t.Errorf("Expected error type to be '%s', was '%s'", playerErrType, err.Type)
	}
}

func GetFakePlayer(name string) *arena.Player {
	username := fmt.Sprintf("%s username", name)
	return &arena.Player{
		Id:        1,
		Name:      name,
		Username:  username,
		InviteUrl: "http://example.com/invite",
		MatchUrl:  "http://example.com/match",
	}
}

// Sending an invitation without specifying a game should return a 400 Bad
// Request.
func TestInviteInvalidMove(t *testing.T) {
	r := mux.NewRouter()
	buf := bytes.NewBufferString("{\"Game\": \"fourup\", \"FirstMove\": \"invalid-parameter\"}")

	// Reassign the player getter function to return no database rows.
	getDatastore = func() engine.Datastore {
		return engine.DummyDatastore{}
	}
	defer reassignDatastoreGetter()

	r.HandleFunc("/players/{player}/invitations", InvitationsHandler)
	req, _ := http.NewRequest("POST", "http://localhost/players/kevinburke/invitations", buf)
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != 400 {
		t.Errorf("Expected status 400 but got %d", resp.Code)
	}
	var err Error
	decodingErr := json.Unmarshal(resp.Body.Bytes(), &err)
	if decodingErr != nil {
		fmt.Println(string(resp.Body.Bytes()))
		t.Fatalf(decodingErr.Error())
	}
	errMsg := "First move value was invalid-parameter but that player is not playing this game"
	if err.Message != errMsg {
		t.Errorf("Expected error message to read '%s', was '%s'", errMsg, err.Message)
	}
	if err.Type != "invalid-first-move" {
		t.Errorf("Expected error type to be 'invalid-first-move', was '%s'", err.Type)
	}
}

//type MatchResponses struct {
//Matches []arena.MatchResponse `json:"matches"`
//}

//func TestOneMatch(t *testing.T) {
//r := mux.NewRouter()
//r.HandleFunc("/games/four-up/matches", MatchesHandler)
//req, _ := http.NewRequest("GET", "http://localhost/games/four-up/matches", nil)

//match := &arena.FourUpMatch{
//Id:     3,
//Winner: nil,
//}
//matchesGetter = func() ([]*arena.FourUpMatch, error) {
//return []*arena.FourUpMatch{match}, nil
//}

//defer reassignMatchesGetter(arena.GetMatches)

//resp := httptest.NewRecorder()
//r.ServeHTTP(resp, req)

//var response MatchResponses
//bits := resp.Body.Bytes()
//fmt.Println(string(bits))
//err := json.Unmarshal(bits, &response)
//if err != nil {
//t.Fatalf(err.Error())
//}
//if len(response.Matches) != 1 {
//t.Fatalf("match length should have been 0, was %d", len(response.Matches))
//}
//mr := response.Matches[0]
//fmt.Println(mr.Winner)
//if mr.Id != 3 {
//t.Errorf("id should have been 3, was %d", mr.Id)
//}
//}

//func reassignMatchGetter(to func(int) (*arena.FourUpMatch, error)) {
//matchGetter = to
//}

//func reassignMatchesGetter(to func() ([]*arena.FourUpMatch, error)) {
//matchesGetter = to
//}

func TestSendInviteOK(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("expected method to be POST, instead was %s", r.Method)
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
