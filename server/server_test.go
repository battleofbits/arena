package main

import (
	"encoding/json"
	//"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
//func setup() {
//// test server
//mux = http.NewServeMux()
//server = httptest.NewServer(mux)

//// github client configured to use test server
//url, _ := url.Parse(server.URL)
//client.BaseURL = url
//client.UploadURL = url
//}

// teardown closes the test HTTP server.
//func teardown() {
//server.Close()
//}

// just testing the infrastructure
func TestHelloWorld(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(HelloWorldHandler))
	defer testServer.Close()
	res, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if string(body) != "Hello World" {
		t.Errorf("Expected 'Hello World' but got %s", string(body))
	}
}

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

//func buildTestRequest(method string, path string, body string, headers map[string][]string, cookies []*http.Cookie) *http.Request {
//host := "127.0.0.1"
//port := "80"
//rawurl := "http://" + host + ":" + port + path
//url_, _ := url.Parse(rawurl)
//proto := "HTTP/1.1"

//if headers == nil {
//headers = map[string][]string{}
//}

//headers["User-Agent"] = []string{"web.go test"}
//if method == "POST" {
//headers["Content-Length"] = []string{fmt.Sprintf("%d", len(body))}
//if headers["Content-Type"] == nil {
//headers["Content-Type"] = []string{"text/plain"}
//}
//}

//req := http.Request{Method: method,
//URL:    url_,
//Proto:  proto,
//Host:   host,
//Header: http.Header(headers),
//Body:   ioutil.NopCloser(bytes.NewBufferString(body)),
//}

//for _, cookie := range cookies {
//req.AddCookie(cookie)
//}
//return &req
//}

//func TestMovesList(t *testing.T) {
//// XXX Need some way to reset the database, or similar, here, so you can
//// actually test interesting things about the list response
////go doServer()
//resp, err := http.Get("http://0.0.0.0:9999/games/four-up/matches/1/moves")
//checkError(err)
//defer resp.Body.Close()
//body, err := ioutil.ReadAll(resp.Body)
//fmt.Println(string(body))
//t.Errorf("XXX: this fails the moves test so you can see the console output")
//}

//func TestMovesList(t *testing.T) {
//// XXX Need some way to reset the database, or similar, here, so you can
//// actually test interesting things about the list response
//go doServer()
//resp, err := http.Get("http://0.0.0.0:9999/games/four-up/matches/1/moves")
//checkError(err)
//defer resp.Body.Close()
//body, err := ioutil.ReadAll(resp.Body)
//fmt.Println(string(body))
//t.Errorf("XXX: this fails the moves test so you can see the console output")
//}

//func TestPlayersList(t *testing.T) {
//// XXX Need some way to reset the database, or similar, here, so you can
//// actually test interesting things about the list response
////go doServer()
//resp, err := http.Get("http://0.0.0.0:9999/players")
//checkError(err)
//defer resp.Body.Close()
//body, err := ioutil.ReadAll(resp.Body)
//fmt.Println(string(body))
//t.Errorf("XXX: this fails the test so you can see the console output")
//}
