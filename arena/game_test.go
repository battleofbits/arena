package main

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
