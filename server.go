package arena

import (
	"encoding/json"
	"github.com/hoisie/web"
	"time"
)

func players(ctx *web.Context) []byte {
	ctx.SetHeader("Content-Type", "application/json", true)
	jsonPlayers, err := json.Marshal(GetPlayers())
	checkError(err)
	return jsonPlayers
}

type Move struct {
	Player string    `json:"player"`
	Column int       `json:"column"`
	Played time.Time `json:"played"`
}

type Moves struct {
	Moves []*Move `json:"moves"`
}

func moves(ctx *web.Context, matchId string) []byte {
	db := getConnection()
	// XXX do a join here to get player name
	query := "SELECT fourup_column, player, played FROM fourup_moves WHERE id = $1"
	rows, err := db.Query(query, matchId)
	checkError(err)
	var moves []*Move
	for rows.Next() {
		var m Move
		var pId int
		err = rows.Scan(&m.Column, &pId, &m.Played)
		checkError(err)
		player, err := GetPlayerById(pId)
		checkError(err)
		player.SetHref()
		m.Player = player.Href
		moves = append(moves, &m)
	}
	jsonMoves, err := json.Marshal(moves)
	checkError(err)
	return jsonMoves
}

func doServer() {
	web.Get("/players", players)
	web.Get("/games/four-up/matches/([^/]+)/moves", moves)
	web.Run("0.0.0.0:9999")
}
