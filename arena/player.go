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
	Username  string `json:"username"`
	MatchUrl  string `json:"-"`
	InviteUrl string `json:"-"`
	Href      string `json:"href"`
}

type Players struct {
	Players []*Player `json:"players"`
}

// Set the href
func (p *Player) SetHref() {
	p.Href = fmt.Sprintf("https://battleofbits.com/players/%s", p.Name)
}

func GetPlayers() ([]*Player, error) {
	db := GetConnection()
	defer db.Close()
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

func CreatePlayer(username string, name string, matchUrl string,
	inviteUrl string) (*Player, error) {
	db := GetConnection()
	defer db.Close()
	player := &Player{
		Username:  username,
		Name:      name,
		MatchUrl:  matchUrl,
		InviteUrl: inviteUrl,
	}
	query := "INSERT INTO players " +
		"(username, name, match_url, invite_url)" +
		"VALUES ($1, $2, $3, $4) RETURNING id"
	err := db.QueryRow(query, username, name, matchUrl, inviteUrl).Scan(
		&player.Id)
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
	return GetPlayer("name", name)
}

func GetPlayerById(playerId int) (*Player, error) {
	return GetPlayer("id", playerId)
}

func GetPlayer(attr string, value interface{}) (*Player, error) {
	var p Player
	db := GetConnection()
	defer db.Close()
	query := fmt.Sprintf("SELECT id, username, name, match_url, invite_url "+
		"FROM players WHERE %s = $1", attr)
	err := db.QueryRow(query, value).Scan(&p.Id, &p.Username, &p.Name,
		&p.MatchUrl, &p.InviteUrl)
	if err != nil {
		return &Player{}, err
	} else {
		p.SetHref()
		return &p, nil
	}
}
