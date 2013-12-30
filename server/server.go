package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/battleofbits/arena/arena"
	"github.com/gorilla/mux"
	"github.com/hoisie/web"
	"log"
	"net/http"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

type Move struct {
	Player string    `json:"player"`
	Column int       `json:"column"`
	Played time.Time `json:"played"`
}

type Moves struct {
	Moves []*Move `json:"moves"`
}

type InviteRequest struct {
	Game             string `json:"game"`
	RequestingPlayer string `json:"requesting_player"`
	FirstMove        string `json:"first_move"`
}

// This is reassigned in tests
var moveGetter = getMoves

func getMoves(moveId int) []*Move {
	db := arena.GetConnection()
	// XXX do a join here to get player name
	query := "SELECT fourup_column, player, played FROM fourup_moves WHERE match_id = $1"
	rows, err := db.Query(query, moveId)
	checkError(err)
	var moves []*Move
	for rows.Next() {
		var m Move
		var pId int
		err = rows.Scan(&m.Column, &pId, &m.Played)
		checkError(err)
		player, err := arena.GetPlayerById(pId)
		checkError(err)
		player.SetHref()
		m.Player = player.Href
		moves = append(moves, &m)
	}
	return moves
}

type MatchResponse struct {
	Id          int64                                  `json:"id"`
	CurrentMove string                                 `json:"current_move"`
	Winner      *NullableString                        `json:"winner"`
	RedPlayer   string                                 `json:"red_player"`
	BlackPlayer string                                 `json:"black_player"`
	Started     time.Time                              `json:"started"`
	Board       *[arena.NumRows][arena.NumColumns]int8 `json:"board"`
	Finished    *time.Time                             `json:"finished"`
}

// Sends an invitation to the invite URL, waits for a response, parses it, etc.
func SendInvite(inviteUrl string, game string, firstMove string) error {
	inviteStruct := &InviteRequest{
		Game: game,
		// XXX, do authentication or use a URL parameter
		RequestingPlayer: "kevinburke",
		FirstMove:        firstMove,
	}
	inviteBody, err := json.Marshal(inviteStruct)
	checkError(err)
	res, err := arena.MakeRequest(inviteUrl, inviteBody)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	if 200 <= res.StatusCode && res.StatusCode <= 299 {
		return nil
	} else {
		return errors.New(fmt.Sprintf(
			"Received error status %s from invite url %s",
			res.Status, inviteUrl,
		))
	}
}

type Response map[string]interface{}

type NullableString struct {
	stringValue string
	isNil       bool
}

func (n NullableString) MarshalJSON() ([]byte, error) {
	if n.isNil == true {
		return []byte{}, nil
	} else {
		return []byte(n.stringValue), nil
	}
}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func moves(ctx *web.Context, matchId string) []byte {
	db := arena.GetConnection()
	// XXX do a join here to get player name
	query := "SELECT fourup_column, player, played FROM fourup_moves WHERE match_id = $1"
	rows, err := db.Query(query, matchId)
	checkError(err)
	var moves []*Move
	for rows.Next() {
		var m Move
		var pId int
		err = rows.Scan(&m.Column, &pId, &m.Played)
		checkError(err)
		player, err := arena.GetPlayerById(pId)
		checkError(err)
		player.SetHref()
		m.Player = player.Href
		moves = append(moves, &m)
	}
	jsonMoves, err := json.Marshal(moves)
	checkError(err)
	return jsonMoves
}

// Common middleware for all API functions
var headerMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Print(req.Method, " ", req.URL)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func DoServer() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/players", headerMiddleware(PlayersHandler)).Methods("GET")
	r.Handle("/players/{player}", headerMiddleware(PlayerHandler)).Methods("GET")
	r.Handle("/games/four-up/matches/{match}/moves",
		headerMiddleware(MovesHandler)).Methods("GET")
	r.Handle("/players/{player}/invitations",
		headerMiddleware(InvitationsHandler)).Methods("POST")
	http.Handle("/", r)
	return r
}

func main() {
	router := DoServer()
	log.Fatal(http.ListenAndServe(":8080", router))
}
