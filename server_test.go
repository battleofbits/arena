package arena

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPlayersList(t *testing.T) {
	// XXX Need some way to reset the database, or similar, here, so you can
	// actually test interesting things about the list response
	go doServer()
	resp, err := http.Get("http://0.0.0.0:9999/players")
	checkError(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	t.Errorf("test failed!")
}
