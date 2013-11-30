package arena

import (
	"encoding/json"
	//"fmt"
	"github.com/hoisie/web"
)

func players(ctx *web.Context) []byte {
	ctx.SetHeader("Content-Type", "application/json", true)
	jsonPlayers, err := json.Marshal(GetPlayers())
	checkError(err)
	return jsonPlayers
}

func doServer() {
	web.Get("/players", players)
	web.Run("0.0.0.0:9999")
}
