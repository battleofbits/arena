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
	Name string `json:"name"`
	// The player's friendly name
	Username string `json:"username"`
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

func GetPlayers() ([]*Player, error) {
	db := getConnection()
	rows, err := db.Query("SELECT username, name from players")
	if err != nil {
		return nil, err
	}
	var players []*Player
	for rows.Next() {
		var p Player
		err = rows.Scan(&p.Username, &p.Name)
		if err != nil {
			return nil, err
		}
		p.SetHref()
		players = append(players, &p)
	}
	return players, nil
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

func GetPlayerById(playerId int) (*Player, error) {
	var p Player
	db := getConnection()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM players WHERE id = $1", playerId).Scan(&p.Id, &p.Username, &p.Name, &p.Url)
	if err != nil {
		return &Player{}, err
	} else {
		return &p, nil
	}
}
