package main

import (
	"database/sql"
	"fmt"
	"github.com/battleofbits/arena/arena"
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
	player, err := arena.GetPlayerByName(playerName)
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
	fmt.Fprint(w, Response{"players": []*arena.Player{player}})
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
		// XXX check ordering here
		mtch, err := arena.CreateMatch(invitedPlayer, requestingPlayer)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, Response{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		mr := &MatchResponse{
			Id:          mtch.Id,
			CurrentMove: mtch.CurrentPlayer.Name,
			Winner:      nil,
			Started:     mtch.Started,
			Finished:    nil,
			Board:       mtch.Board,
			RedPlayer:   mtch.RedPlayer.Name,
			BlackPlayer: mtch.BlackPlayer.Name,
		}
		fmt.Fprint(w, Response{"matches": []*MatchResponse{mr}})
		// XXX check ordering here
		arena.StartMatch(mtch, invitedPlayer, requestingPlayer)
	}
})

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
