package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	USER_AGENT = "battleofbits/0.2"
	READ_TIME  = 30
)

type Player struct {
	Id int64 `json:"-"`
	// The player's friendly name
	Name      string `json:"name"`
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

// Every instance of a game must implement this interface
type Match interface {

	// Retrieve the player whose turn it is.
	CurrentPlayer() *Player

	// Serialize a move from the byte string and apply it to the board. Returns
	// `true` if the game is over, and an error if the move was unreadable or
	// invalid.
	Play(*Player, []byte) (bool, error)

	// Retrieve the winning player.
	Winner() Player

	// Determine whether the match is a stalemate.
	Stalemate() bool

	// Advance the turn.
	NextPlayer() *Player
}

// Make a request to a player's URL
func MakeRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", USER_AGENT)

	// XXX, make this anonymous.
	type HttpTimeoutResponse struct {
		Resp *http.Response
		Err  error
	}

	httpRes := make(chan HttpTimeoutResponse, 1)
	go func() {
		res, err := client.Do(req)
		httpRes <- HttpTimeoutResponse{Resp: res, Err: err}
	}()
	select {
	case res := <-httpRes:
		return res.Resp, res.Err
	case <-time.After(time.Second * READ_TIME):
		return nil, errors.New("HTTP Response was not received in time")
	}
}

// Assemble and make an HTTP request to the user's URL
// Returns the column of the response
func GetMove(match Match, player Player) ([]byte, error) {
	payload, err := json.Marshal(match)

	if err != nil {
		return []byte{}, err
	}

	response, err := MakeRequest(player.MatchUrl, payload)

	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func PlayMatch(match Match, datastore Datastore) error {
	if datastore == nil {
		datastore = PostgresDatastore{
			url: "arena@localhost:5432/arena?sslmode=disable",
		}
	}
	datastore.SerializeMatch(match)
	for {
		player := match.CurrentPlayer()

		move, err := GetMove(match, *player)

		if err != nil {
			return fmt.Errorf("Player's server %s was unreachable: %s", player.MatchUrl, err)
		}

		gameover, err := match.Play(player, move)
		datastore.SerializeMatch(match)

		player = match.NextPlayer()

		if err != nil {
			//Move was invalid, game is over
			return err
		}

		if gameover {
			if match.Stalemate() {

			}
			// Record the winner
			_ = match.Winner()
			return nil
		}
	}
}
