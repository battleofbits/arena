package engine

import (
	"fmt"
	"database/sql"
	"github.com/lib/pq"
)

type Datastore interface {
	CreatePlayer(string, string, string, string) (Player, error)
	GetPlayers() []Player
	GetPlayerByName(string) (Player, error)
	GetPlayerById(int) (Player, error)
}

// rawUrl := "postgres://postgres_arena@localhost:5432/arena?sslmode=disable"
func (p Postgres) getConnection() (*sql.DB, error) {
	url, err := pq.ParseURL(p.url)

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return db, nil
}

type Postgres struct {
	url string
}


func (p Postgres) GetPlayers() ([]*Player, error) {
	db, err := p.getConnection()
	if err != nil {
		return nil, err
	}
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

func (p Postgres) CreatePlayer(username string, name string, matchUrl string) (*Player, error) {
	db, err := p.getConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	player := &Player{
		Username:  username,
		Name:      name,
		MatchUrl:  matchUrl,
	}
	query := "INSERT INTO players " +
		"(username, name, match_url, invite_url)" +
		"VALUES ($1, $2, $3, $4) RETURNING id"
	err = db.QueryRow(query, username, name, matchUrl, "").Scan(
		&player.Id)
	var pqerr *pq.Error
	if err != nil {
		pqerr = err.(*pq.Error)
	}
	if pqerr != nil && pqerr.Code.Name() == "unique_violation" {
		return p.GetPlayerByName(name)
	}
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (p Postgres) GetPlayerByName(name string) (*Player, error) {
	return p.getPlayer("name", name)
}

func (p Postgres) GetPlayerById(playerId int) (*Player, error) {
	return p.getPlayer("id", playerId)
}

func (pd Postgres) getPlayer(attr string, value interface{}) (*Player, error) {
	var p Player
	db, err := pd.getConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT id, username, name, match_url, invite_url "+
		"FROM players WHERE %s = $1", attr)
	err = db.QueryRow(query, value).Scan(&p.Id, &p.Username, &p.Name,
		&p.MatchUrl, "")
	if err != nil {
		return &Player{}, err
	} else {
		p.SetHref()
		return &p, nil
	}
}
