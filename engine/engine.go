package engine

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const USER_AGENT = "battleofbits/0.1"

// The player's friendly name
type Player struct {
	Id       int64  `json:"-"`
	Name     string `json:"name"`
	Username string `json:"username"`
	MatchUrl string `json:"-"`
	Href     string `json:"href"`
}

type Match interface {
	CurrentPlayer() Player
	Play(Player, []byte) (bool, error)
	Winner() (Player, error)
	Stalemate() bool
}

func MakeRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(postBody))

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", USER_AGENT)
	// XXX, set a timeout here
	return client.Do(req)
}

// Assemble and make an HTTP request to the user's URL
// Returns the column of the response
func GetMove(match *Match) ([]byte, error) {
	payload, err := json.Marshal(match)

	if err != nil {
		return []byte{}, err
	}

	player := match.CurrentPlayer()

	response, err := MakeRequest(player.MatchUrl, payload)

	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return []byte{}, err
	}

	return body
}

func PlayMatch(match *Match, players []Player) {
	for {
		player := m.CurrentPlayer()

		move, err := GetMove(player, m)

		if err != nil {
			break
		}

		gameover, err := m.Play(player, move)

		if err != nil {
			//Move was invalid, game is over
			break
		}

		if gameover {
			// Record the winner
			winner, err := m.Winner()
		}
	}
}
