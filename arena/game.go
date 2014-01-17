// Contains the base logic for a game, integrates all of the models etc.
package arena

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	//"net/http"
)

const URL = "http://localhost:5000/fourup"
const BaseUri = "https://battleofbits.com"

const NumRows = 6
const NumColumns = 7
const NumConsecutive = 4

//type TurnPlayers struct {
//Red   string `json:"R"`
//Black string `json:"B"`
//}

//type FourUpTurn struct {
//Href     string                      `json:"href"`
//Players  *TurnPlayers                `json:"players"`
//Turn     string                      `json:"turn"`
//Loser    string                      `json:"loser"`
//Winner   string                      `json:"winner"`
//Started  string                      `json:"started"`
//Finished string                      `json:"finished"`
//Moves    string                      `json:"moves"`
//Board    [NumRows][NumColumns]string `json:"board"`
//}

//func DoForfeit(loser *Player, reason error) {
//fmt.Println(fmt.Sprintf("player %s forfeits because of %s", loser.Username, reason.Error()))
//}

//func DoGameOver(match *FourUpMatch, winner *Player, loser *Player) {
//fmt.Println("Game is over. Winner is ", winner.Username, ". Notifying winner and loser...")
//MarkWinner(match, winner)
//NotifyWinner(winner)
//NotifyLoser(loser)
//}

// Assemble and make an HTTP request to the user's URL
// Returns the column of the response

func NotifyWinner(winner *Player) {
	fmt.Println("Notifying winner...")
}

func NotifyLoser(loser *Player) {
	fmt.Println("Notifying loser...")
}

// In the background, let people know about the new move
//func NotifySubscribers(move int8, match *FourUpMatch) {

//}

// playerId - 1 for red, 2 for black. XXX, refactor this.
//func DoPlayerMove(player *Player, otherPlayer *Player, match *FourUpMatch, playerId int) error {
//move, err := GetMove(match)
//if err != nil {
//DoForfeit(player, err)
//DoGameOver(match, otherPlayer, player)
//return err
//}
//err = DoNewMove(move, match)
//if err != nil {
//// XXX, do game over here, or switch based on the error type, etc.
//return err
//}
//checkError(err)
//if GameOver(*match.Board) {
//DoGameOver(match, player, otherPlayer)
//return errors.New("Game is over.")
//}
//if IsBoardFull(*match.Board) {
//DoTieGame(match, player, otherPlayer)
//return err
//}
//return nil
//}

//func DoMatch(match *FourUpMatch, redPlayer *Player, blackPlayer *Player) *FourUpMatch {
//for {
//match.CurrentPlayer = redPlayer
//err := DoPlayerMove(redPlayer, blackPlayer, match, 1)
//// XXX, evaluate positioning of this update.
//if err != nil {
//break
//}

//match.CurrentPlayer = blackPlayer
//err = DoPlayerMove(blackPlayer, redPlayer, match, 2)
//// XXX, evaluate logic here
//if err != nil {
//break
//}
//}
//return match
//}
