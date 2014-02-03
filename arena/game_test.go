package arena

import (
	"fmt"
	"testing"
)

func TestGetMatchHref(t *testing.T) {
	t.Parallel()
	url := "https://battleofbits.com/games/four-up/matches/1"
	if getMatchHref(1) != url {
		t.Errorf(fmt.Sprintf("Expected url %s but got %s", url, getMatchHref(1)))
	}
}

//func TestSerializeTurn(t *testing.T) {
//t.Parallel()
//red := &Player{
//Name:     "kevin",
//Username: "kb",
//}
//black := &Player{
//Name:     "deferman",
//Username: "kyle who cant spell",
//}

//match := CreateFourUpMatch(red, black)
//turn := serializeTurn(match)
//if turn.Turn != BaseUri+"/players/kevin" {
//t.Errorf("it should be kevin's turn but instead was", turn.Turn)
//}
//if turn.Players.Black != BaseUri+"/players/deferman" {
//t.Errorf("black player should be deferman but was", turn.Players.Black)
//}
//}

//func TestDoGame(t *testing.T) {
//redPlayer, redErr := CreatePlayer("Kevin Burke", "kevinburke", URL)
//if redErr != nil {
//t.Fatalf(redErr.Error())
//}
//blackPlayer, _ := CreatePlayer("Kyle Conroy", "kyleconroy", URL)
//match := CreateFourUpMatch(redPlayer, blackPlayer)
//dbErr := WriteMatch(match)
//if dbErr != nil {
//t.Fatalf(dbErr.Error())
//}
//match = DoMatch(match, redPlayer, blackPlayer)
//}
