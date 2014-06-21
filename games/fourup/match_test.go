package fourup

import (
	"fmt"
	"github.com/battleofbits/arena/engine"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGameLogic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"column\": %d}", rand.Intn(NumColumns))
	}))

	players := []*engine.Player{
		&engine.Player{Id: 1, Name: "kevinburke", MatchUrl: ts.URL},
		&engine.Player{Id: 2, Name: "deferman", MatchUrl: ts.URL},
	}
	match, err := CreateMatch(players)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = engine.PlayMatch(match, engine.DummyDatastore{})
	if err != nil {
		t.Errorf(err.Error())
	}
	// XXX, make this a better test.
	fmt.Println(match.Board)
}
