package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/battleofbits/arena/arena"
	"github.com/battleofbits/arena/engine"
	"github.com/battleofbits/arena/games/fourup"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"strconv"
)

var PlayersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	players, err := arena.GetPlayers()
	checkError(err)
	fmt.Fprint(w, Response{"players": players})
})

var PlayerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	playerName := mux.Vars(r)["player"]
	datastore := getDatastore()
	player, err := datastore.GetPlayerByName(playerName)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, Response{
				"error": fmt.Sprintf("No players with name %s", playerName),
			})
		} else {
			// XXX, middleware etc
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, Response{"error": err.Error()})
		}
		return
	}
	fmt.Fprint(w, Response{"players": []*engine.Player{&player}})
})

//var matchesGetter = arena.GetMatches
//var matchGetter = arena.GetMatch

//var MatchHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//matchId := mux.Vars(r)["match"]
//id, err := strconv.Atoi(matchId)
//if err != nil {
//w.WriteHeader(http.StatusBadRequest)
//fmt.Fprint(w, Response{"error": err.Error()})
//return
//}
//match, err := matchGetter(id)
//if err != nil {
//if err == sql.ErrNoRows {
//w.WriteHeader(http.StatusNotFound)
//fmt.Fprint(w, Response{
//"error": fmt.Sprintf("No matches with id %s", matchId),
//})
//return
//} else {
//w.WriteHeader(http.StatusBadRequest)
//fmt.Fprint(w, Response{"error": err.Error()})
//return
//}
//}
//fmt.Fprint(w, Response{"matches": []*arena.FourUpMatch{match}})
//})

//var MatchesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//matches, err := matchesGetter()
//if err != nil {
//w.WriteHeader(http.StatusInternalServerError)
//fmt.Println(err.Error())
//fmt.Fprint(w, Response{"error": "We experienced an error. Please try again"})
//return
//}
//if len(matches) == 0 {
//// json.Marshal returns null instead of an empty list for a pointer
//// with no data.
//fmt.Fprint(w, Response{"matches": []string{}})
//} else {
//fmt.Fprint(w, Response{"matches": matches})
//}
//})

type InviteBody struct {
	Game      string
	Player    string
	FirstMove string
}

func abortWithError(err error, w http.ResponseWriter) {
	abortWithTypedError(err, "", w)
}

func abortWithTypedError(err error, typ string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, Response{"type": typ, "message": err.Error()})
}

var getDatastore = func() engine.Datastore {
	return engine.GetPostgresDatastore()
}

var originalDatastoreGetter = getDatastore

var reassignDatastoreGetter = func() {
	getDatastore = originalDatastoreGetter
}

// Handle an invitation to play a new game
var InvitationsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	datastore := getDatastore()
	var ivb InviteBody
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		abortWithError(err, w)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &ivb)
	if err != nil {
		abortWithError(err, w)
		return
	}

	if err != nil {
		abortWithError(err, w)
		return
	}
	if ivb.Game == "" {
		abortWithTypedError(errors.New("No game specified"), "invalid-game", w)
		return
	}
	invitedPlayerName := mux.Vars(r)["player"]
	if invitedPlayerName == "" {
		abortWithTypedError(errors.New("No player specified"), "invalid-player", w)
		return
	}

	invitedPlayer, err := datastore.GetPlayerByName(invitedPlayerName)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, Response{
				"message": fmt.Sprintf("No players with name %s", invitedPlayerName),
				"type":    "invalid-player",
			})
			return
		} else {
			abortWithError(err, w)
			return
		}
		return
	}

	// XXX, pull from authentication or parameters
	requestingPlayerName := "kevinburke"

	// XXX modularize this and above
	requestingPlayer, err := datastore.GetPlayerByName(requestingPlayerName)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, Response{
				"message": fmt.Sprintf("No players with name %s", requestingPlayerName),
				"type":    "invalid-requesting-player",
			})
		} else {
			abortWithError(err, w)
		}
		return
	}

	var playerWithFirstMove string
	if ivb.FirstMove == "random" || ivb.FirstMove == "" {
		if rand.Intn(2) == 0 {
			playerWithFirstMove = requestingPlayerName
		} else {
			playerWithFirstMove = invitedPlayerName
		}
	} else if ivb.FirstMove != requestingPlayerName &&
		ivb.FirstMove != invitedPlayerName {
		msg := fmt.Sprintf("First move value was %s but that player is not playing this game", ivb.FirstMove)
		abortWithTypedError(errors.New(msg), "invalid-first-move", w)
		return
	} else {
		playerWithFirstMove = ivb.FirstMove
	}

	err = SendInvite(invitedPlayer.InviteUrl, ivb.Game, playerWithFirstMove)
	if err != nil {
		abortWithError(err, w)
		return
	}

	// Fork off a goroutine to run the match.

	//startNullable := &arena.NullTime{
	//Valid: true,
	//Time:  mtch.Started,
	//}
	//finishedNullable := &arena.NullTime{
	//Valid: false,
	//}
	//mr, err := json.Marshal(mtch)
	////w.Header().Set("Location", mr.Href)
	//w.WriteHeader(http.StatusCreated)
	//fmt.Fprint(w, Response{"matches": mr})

	players := []*engine.Player{&invitedPlayer, &requestingPlayer}
	match, err := fourup.CreateMatch(players)
	if err != nil {
		abortWithError(err, w)
	}
	engine.PlayMatch(match, datastore)
	//}
})

// This is reassigned in tests
var moveGetter = getMoves

var MovesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["match"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, Response{"error": err.Error()})
		return
	}
	moves := moveGetter(id)
	fmt.Fprint(w, Response{"moves": moves})
})
