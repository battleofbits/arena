package arena

import (
	"fmt"
	"testing"
)

func TestGetHref(t *testing.T) {
	url := "https://battleofbits.com/games/four-up/matches/1"
	if getHref(1) != url {
		t.Errorf(fmt.Sprintf("Expected url %s but got %s", url, getHref(1)))
	}
}

func TestDoGame(t *testing.T) {
	redPlayer, _ := CreatePlayer("Kevin Burke", "kevinburke", URL)
	blackPlayer, _ := CreatePlayer("Kyle Conroy", "kyleconroy", URL)
	match, fourupErr := CreateFourUpMatch(redPlayer, blackPlayer)
	checkError(fourupErr)
	match = DoMatch(match, redPlayer, blackPlayer)
}
