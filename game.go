// Contains the base logic for a game, integrates all of the models etc.

// Webhook API
//
// Each player has an webhook associated webhook for a specifc game. For example,
// the player `deepblue` has the following URL associated for Four Up.
//
//	http://example.com/fourup/hook
//
// When it's deepblue's turn, we'll POST to that endpoint with a board state. The
// endpoint has 30 seconds to respond with a valid move. If the endpoint takes
// longer than 30 seconds, or the returned move isn't valid, the game is
// forfeited.  So don't make an invalid move!
//
// Example Webhook
//
// 	POST /fourup/hook
// 	Content-Type: application/json
//
// 	{
// 	  "href": "https://battleofbits.com/games/four-up/matches/1",
// 	  "players": {
// 	    "https://battleofbits.com/players/deepblue": 1,
// 	    "https://battleofbits.com/players/garry": 2
// 	  },
// 	  "turn": "https://battleofbits.com/players/deepblue",
// 	  "loser": "",
// 	  "winner": "",
// 	  "started": "2013-01-01T23:00:01Z",
// 	  "finished": "",
// 	  "moves": "https://battleofbits.com/games/four-up/matches/1/moves",
// 	  "board": [
//	    [0, 0, 0, 0, 0, 0, 0],
//	    [0, 0, 0, 0, 0, 0, 0],
//	    [0, 0, 0, 0, 0, 0, 0],
//	    [0, 0, 0, 0, 0, 0, 0],
//	    [0, 1, 0, 0, 0, 0, 0],
//	    [0, 2, 1, 0, 0, 0, 0]
// 	  ]
// 	}
//
// A valid response is a JSON move.
//
// 	HTTP/1.0 200 OK
// 	Content-Type: application/json
//
// 	{
// 	  "column": 2
// 	}
//
// ## API
//
// Each board game will have a custom media type associated with it, so that
// others can utilize the JSON representations, even outside of battleofbits.
//
// ### GET /players
//
// ```js
// {
//   "players": [{
//     "href": "https://battleofbits.com/players/deepblue",
//     "username": "deepblue",
//     "name": "Deep Blue"
//   },{
//     "href": "https://battleofbits.com/players/garry",
//     "username": "garry",
//     "name": "Garry Kasparov"
//   }]
// }
// ```
//
// ### POST /players
//
// Create a new player. Username is required
//
// ### GET /games
//
// ```js
// HTTP/1.1 200 OK
// Content-Type: application/json
//
// {
//   "games": [{
//     "name": "Four Up",
//     "href": "https://battleofbits.com/games/four-up",
//     "matches": "https://battleofbits.com/games/four-up/matches"
//   }]
// }
// ```
//
// ### GET /games/four-up/matches
//
// Return a list of all present and on-going matches
//
// ```js
// HTTP/1.1 200 OK
// Content-Type: application/json
//
// {
//   "matches": [{
//     "href": "https://battleofbits.com/games/four-up/matches/1",
//     "players": {
//       "R": "https://battleofbits.com/players/deepblue",
//       "B": "https://battleofbits.com/players/garry"
//     },
//     "winner": "",
//     "started": "2013-01-01T23:00:01Z",
//     "finished": "",
//     "moves": "https://battleofbits.com/games/four-up/matches/1/moves",
//     "board": [
//       ["","","","","","",""],
//       ["","","","","","",""],
//       ["","","","","","",""],
//       ["","","","","","",""],
//       ["","B","","","","",""],
//       ["","R","B","","","",""]
//     ]
//   }]
// }
// ```
//
// ### GET /games/four-up/matches/{id}/moves
//
// ```js
// {
//   "moves": [{
//     "player": "https://battleofbits.com/players/deepblue",
//     "column": 2,
//     "played": "2013-01-01T23:00:01Z",
//   }]
// }
// ```
//
// ### SUBSCRIBE /games/four-up/matches/{id}/moves
//
// Pass the correct accept header
//
// ```
// Accept: text/event-stream
// ```
//
// ```
// HTTP/1.1 200 OK
// Content-Type: text/event-stream
//
// event: patch
// data: [
//         {
//           "op": "replace",
//           "path": "/games/four-up/matches/{id}/board"
//           "value": [
// 	      ["","","","","","",""],
// 	      ["","","","","","",""],
// 	      ["","","","","","",""],
// 	      ["","","","","","",""],
// 	      ["","","","","","",""],
// 	      ["","","","","","",""]
// 	    ]
//         }
//       ]

package fourup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const URL = "http://localhost:5000/fourup"
const BaseUri = "https://battleofbits.com"

type TurnPlayers struct {
	Red   string `json:"R"`
	Black string `json:"B"`
}

type FourUpTurn struct {
	Href     string                      `json:"href"`
	Players  *TurnPlayers                `json:"players"`
	Turn     string                      `json:"turn"`
	Loser    string                      `json:"loser"`
	Winner   string                      `json:"winner"`
	Started  string                      `json:"started"`
	Finished string                      `json:"finished"`
	Moves    string                      `json:"moves"`
	Board    [NumRows][NumColumns]string `json:"board"`
}

type FourUpResponse struct {
	Column int `json:"column"`
}

func serializeTurn(match *FourUpMatch) *FourUpTurn {
	return &FourUpTurn{
		Href:  getMatchHref(match.Id),
		Board: GetStringBoard(match.Board),
		Turn:  fmt.Sprintf(BaseUri+"/players/%s", match.CurrentPlayer.Name),
		Players: &TurnPlayers{
			Red:   fmt.Sprintf(BaseUri+"/players/%s", match.RedPlayer.Name),
			Black: fmt.Sprintf(BaseUri+"/players/%s", match.BlackPlayer.Name),
		},
	}
}
func DoForfeit(loser *Player, reason error) {
	fmt.Println(fmt.Sprintf("player %s forfeits because of %s", loser.Username, reason.Error()))
}

func DoGameOver(match *FourUpMatch, winner *Player, loser *Player) {
	fmt.Println("Game is over. Winner is ", winner.Username, ". Notifying winner and loser...")
	MarkWinner(match, winner)
	NotifyWinner(winner)
	NotifyLoser(loser)
}

// Assemble and make an HTTP request to the user's URL
// Returns the column of the response
func GetMove(match *FourUpMatch) (int, error) {
	turn := serializeTurn(match)
	postBody, err := json.Marshal(turn)
	checkError(err)
	httpResponse, err := MakeRequest(match.CurrentPlayer.Url, postBody)
	return ParseResponse(httpResponse)
}

// Retrieves the column from the http response
func ParseResponse(response *http.Response) (int, error) {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}
	var fourUpResponse FourUpResponse
	err = json.Unmarshal(body, &fourUpResponse)
	if err != nil {
		return -1, err
	}
	return fourUpResponse.Column, nil
}

func NotifyWinner(winner *Player) {
	fmt.Println("Notifying winner...")
}

func NotifyLoser(loser *Player) {
	fmt.Println("Notifying loser...")
}

func MarkWinner(match *FourUpMatch, winner *Player) error {
	db := getConnection()
	defer db.Close()
	_, err := db.Exec("UPDATE fourup_matches SET winner = $1, "+
		"finished = NOW() at time zone 'utc' WHERE id = $2",
		winner.Id, match.Id)
	return err
}

// Write a new move to the database
func WriteMove(move int, match *FourUpMatch) (int64, error) {
	db := getConnection()
	defer db.Close()
	query := "INSERT INTO fourup_moves (fourup_column, player, move_number, match_id, played)" +
		"VALUES ($1, $2, $3, $4, NOW() at time zone 'utc') RETURNING id"
	var moveId int64
	err := db.QueryRow(query, move, match.CurrentPlayer.Id, match.MoveId, match.Id).Scan(&moveId)
	return moveId, err
}

// Do a whole bunch of stuff associated with new moves
// Error handling is a little tricky because most of the errors would be
// database or other errors.
func DoNewMove(move int, match *FourUpMatch) error {
	var err error
	match.Board, err = ApplyMoveToBoard(move, int(match.CurrentPlayer.Id), match.Board)
	// XXX
	//if err != nil {
	//DoForfeit(player, err)
	//DoGameOver(match, otherPlayer, player)
	//return err
	//}
	// once we know move was valid, update the database
	_, err = WriteMove(move, match)
	checkError(err)
	match.MoveId++
	err = UpdateMatch(match)
	checkError(err)
	NotifySubscribers(move, match)
	return nil
}

// In the background, let people know about the new move
func NotifySubscribers(move int, match *FourUpMatch) {

}

// playerId - 1 for red, 2 for black. XXX, refactor this.
func DoPlayerMove(player *Player, otherPlayer *Player, match *FourUpMatch, playerId int) error {
	move, err := GetMove(match)
	if err != nil {
		DoForfeit(player, err)
		DoGameOver(match, otherPlayer, player)
		return err
	}
	err = DoNewMove(move, match)
	if err != nil {
		// XXX, do game over here, or switch based on the error type, etc.
		return err
	}
	checkError(err)
	if GameOver(*match.Board) {
		DoGameOver(match, player, otherPlayer)
		return errors.New("Game is over.")
	}
	if IsBoardFull(*match.Board) {
		DoTieGame(match, player, otherPlayer)
		return err
	}
	return nil
}

func DoTieGame(match *FourUpMatch, playerOne *Player, playerTwo *Player) {
	fmt.Println("Tie Game!")
}

func DoMatch(match *FourUpMatch, redPlayer *Player, blackPlayer *Player) *FourUpMatch {
	for {
		match.CurrentPlayer = redPlayer
		err := DoPlayerMove(redPlayer, blackPlayer, match, 1)
		// XXX, evaluate positioning of this update.
		if err != nil {
			break
		}

		match.CurrentPlayer = blackPlayer
		err = DoPlayerMove(blackPlayer, redPlayer, match, 2)
		// XXX, evaluate logic here
		if err != nil {
			break
		}
	}
	return match
}
