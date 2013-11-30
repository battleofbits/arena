package arena

// model for a player

import (
	"fmt"
	"github.com/lib/pq"
)

type Player struct {
	// The autoid for the player
	Id int64 `json:"-"`
	// The player's unique Id
	Name string
	// The player's friendly name
	Username string
	Url      string `json:"-"`
	Href     string `json:"href"`
}

type Players struct {
	Players []*Player `json:"players"`
}

// Set the href
func (p *Player) SetHref() {
	p.Href = fmt.Sprintf("https://battleofbits.com/players/%s", p.Name)
}

func GetPlayers() []*Player {
	db := getConnection()
	rows, err := db.Query("SELECT (username, name) from players")
	checkError(err)
	var players []*Player
	for rows.Next() {
		var p *Player
		err = rows.Scan(&p.Username, &p.Name)
		p.SetHref()
		checkError(err)
		players = append(players, p)
	}
	return players
}

func CreatePlayer(username string, name string, url string) (*Player, error) {
	db := getConnection()
	defer db.Close()
	player := &Player{
		Username: username,
		Name:     name,
		Url:      url,
	}
	err := db.QueryRow("INSERT INTO players (username, name, url) VALUES ($1, $2, $3) RETURNING id", username, name, url).Scan(&player.Id)
	var pqerr *pq.Error
	if err != nil {
		pqerr = err.(*pq.Error)
	}
	if pqerr != nil && pqerr.Code.Name() == "unique_violation" {
		return GetPlayerByName(name)
	}
	checkError(err)
	return player, nil
}

func GetPlayerByName(name string) (*Player, error) {
	var p Player
	db := getConnection()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM players WHERE name = $1", name).Scan(&p.Id, &p.Username, &p.Name, &p.Url)
	if err != nil {
		return &Player{}, err
	} else {
		return &p, nil
	}
}
