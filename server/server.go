package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/battleofbits/arena/arena"
	"github.com/gorilla/mux"
	"github.com/hoisie/web"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var moveGetter = getMoves

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

var PlayersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	players, err := arena.GetPlayers()
	checkError(err)
	fmt.Fprint(w, Response{"players": players})
})

var InvitationsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// XXX, middleware etc
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{"error": err.Error()})
		return
	}
	game := r.Form.Get("Game")
	if game == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{"error": "No game specified"})
		return
	}
	invitedPlayerName := mux.Vars(r)["player"]
	if invitedPlayerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{"error": "No player specified"})
		return
	}
	invitedPlayer, err := arena.GetPlayerByName(invitedPlayerName)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, Response{
				"error": fmt.Sprintf("No players with name %s", invitedPlayerName),
			})
		} else {
			// XXX, middleware etc
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, Response{"error": err.Error()})
		}
		return
	}

	if invitedPlayer == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{
			"error": fmt.Sprintf("player %s not found", invitedPlayerName),
		})
		return
	}

	// XXX, pull from authentication or parameters
	requestingPlayerName := "kevinburke"

	// XXX modularize this and above
	requestingPlayer, err := arena.GetPlayerByName(requestingPlayerName)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, Response{
				"error": fmt.Sprintf("No players with name %s", requestingPlayerName),
			})
		} else {
			// XXX, middleware etc
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, Response{"error": err.Error()})
		}
		return
	}

	if requestingPlayer == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{
			"error": fmt.Sprintf("player %s not found", requestingPlayerName),
		})
		return
	}

	incomingFirstMove := r.Form.Get("FirstMove")
	var playerWithFirstMove string
	if incomingFirstMove == "random" || incomingFirstMove == "" {
		if rand.Intn(2) == 0 {
			playerWithFirstMove = requestingPlayerName
		} else {
			playerWithFirstMove = invitedPlayerName
		}
	} else if incomingFirstMove != requestingPlayerName &&
		incomingFirstMove != invitedPlayerName {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{
			"error": fmt.Sprintf("first move value was %s but player %s is "+
				"not in the game", incomingFirstMove, invitedPlayerName),
		})
		return
	} else {
		playerWithFirstMove = incomingFirstMove
	}

	err = SendInvite(invitedPlayer.InviteUrl, game, playerWithFirstMove)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{"error": err.Error()})
		return
	} else {
		fmt.Fprint(w, Response{"invitation": "success"})
		arena.CreateAndDoMatch(invitedPlayer, requestingPlayer)
	}
})

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

var MovesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["match"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// XXX, 400 error.
		fmt.Println("bad id")
	}
	moves := moveGetter(id)
	fmt.Fprint(w, Response{"moves": moves})
})

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
