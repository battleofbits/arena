package arena

import (
	"testing"
)

func TestSetHref(t *testing.T) {
	t.Parallel()
	p := &Player{
		Name: "foobar",
	}
	p.SetHref()
	url := "https://battleofbits.com/players/foobar"
	if p.Href != url {
		t.Errorf("Player href should have been", url, "instead was", p.Href)
	}
}
