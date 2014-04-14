package engine

import (
	"database/sql"
	"github.com/lib/pq"
)

const (
	queryAllPlayers = `SELECT username, name from players`
	queryNewPlayer  = `INSERT INTO players (username, name, match_url, invite_url) 
				VALUES ($1, $2, $3, $4) RETURNING id`
	queryPlayerId = `SELECT id, username, name, match_url, invite_url
				FROM players WHERE id = $1`
	queryPlayerName = `SELECT id, username, name, match_url, invite_url
				FROM players WHERE name = $1`
)

type Datastore interface {
	CreatePlayer(string, string, string) (Player, error)
	GetPlayers() ([]Player, error)
	GetPlayerByName(string) (Player, error)
	GetPlayerById(int) (Player, error)
	SerializeMatch(Match) error
}

type DummyDatastore struct {
}

type NotFoundDatastore struct {
	DummyDatastore
}

func (p NotFoundDatastore) GetPlayerByName(name string) (Player, error) {
	return Player{}, sql.ErrNoRows
}

func (p NotFoundDatastore) GetPlayerById(playerId int) (Player, error) {
	return Player{}, sql.ErrNoRows
}

func (p DummyDatastore) GetPlayers() ([]Player, error) {
	return []Player{}, nil
}

func (p DummyDatastore) CreatePlayer(username string, name string, matchUrl string) (Player, error) {
	return Player{
		Username: username,
		Name:     name,
		MatchUrl: matchUrl,
	}, nil
}

func (p DummyDatastore) SerializeMatch(match Match) error {
	return nil
}

func (p DummyDatastore) GetPlayerByName(name string) (Player, error) {
	return Player{Name: name}, nil
}

func (p DummyDatastore) GetPlayerById(playerId int) (Player, error) {
	return Player{Id: 0}, nil
}

type PostgresDatastore struct {
	url string
}

func GetPostgresDatastore() Datastore {
	return PostgresDatastore{
		url: "arena@localhost:5432/arena?sslmode=disable",
	}
}

// rawUrl := "PostgresDatastore://PostgresDatastore_arena@localhost:5432/arena?sslmode=disable"
func (p PostgresDatastore) getConnection() (*sql.DB, error) {
	url, err := pq.ParseURL(p.url)

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("PostgresDatastore", url)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p PostgresDatastore) GetPlayers() ([]Player, error) {
	db, err := p.getConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(queryAllPlayers)
	if err != nil {
		return nil, err
	}
	var players []Player
	for rows.Next() {
		var p Player
		err = rows.Scan(&p.Username, &p.Name)
		if err != nil {
			return nil, err
		}
		p.SetHref()
		players = append(players, p)
	}
	return players, nil
}

func (p PostgresDatastore) SerializeMatch(match Match) error {
	return nil
}

func (p PostgresDatastore) CreatePlayer(username string, name string, matchUrl string) (Player, error) {
	db, err := p.getConnection()
	if err != nil {
		return Player{}, err
	}
	defer db.Close()

	player := Player{
		Username: username,
		Name:     name,
		MatchUrl: matchUrl,
	}
	err = db.QueryRow(queryNewPlayer, username, name, matchUrl, "").Scan(
		&player.Id)
	var pqerr *pq.Error
	if err != nil {
		pqerr = err.(*pq.Error)
	}
	if pqerr != nil && pqerr.Code.Name() == "unique_violation" {
		return p.GetPlayerByName(name)
	}
	if err != nil {
		return Player{}, err
	}
	return player, nil
}

func (p PostgresDatastore) GetPlayerByName(name string) (Player, error) {
	return p.getPlayer(queryPlayerName, name)
}

func (p PostgresDatastore) GetPlayerById(playerId int) (Player, error) {
	return p.getPlayer(queryPlayerId, playerId)
}

func (pd PostgresDatastore) getPlayer(query string, value interface{}) (Player, error) {
	var p Player
	db, err := pd.getConnection()
	if err != nil {
		return Player{}, err
	}
	defer db.Close()
	err = db.QueryRow(query, value).Scan(&p.Id, &p.Username, &p.Name,
		&p.MatchUrl, "")
	if err != nil {
		return Player{}, err
	} else {
		p.SetHref()
		return p, nil
	}
}
