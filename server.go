package arena

import (
	"encoding/json"
	"github.com/hoisie/web"
)

func players(ctx *web.Context) []byte {
	ctx.SetHeader("Content-Type", "application/json", true)
	jsonPlayers, err := json.Marshal(GetPlayers())
	checkError(err)
	return jsonPlayers
}

func main() {
	web.Get("/players", players)
}
