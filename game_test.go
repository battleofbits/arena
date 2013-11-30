package arena

import (
	"fmt"
	"testing"
)

func TestGetMatchHref(t *testing.T) {
	url := "https://battleofbits.com/games/four-up/matches/1"
	if getMatchHref(1) != url {
		t.Errorf(fmt.Sprintf("Expected url %s but got %s", url, getMatchHref(1)))
	}
}

func TestDoGame(t *testing.T) {
	redPlayer, _ := CreatePlayer("Kevin Burke", "kevinburke", URL)
	blackPlayer, _ := CreatePlayer("Kyle Conroy", "kyleconroy", URL)
	match, fourupErr := CreateFourUpMatch(redPlayer, blackPlayer)
	if fourupErr != nil {
		t.Fatalf(fourupErr.Error())
	}
	match = DoMatch(match, redPlayer, blackPlayer)
}
