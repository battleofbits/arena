package arena

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestNullTimeJSON(t *testing.T) {
	foo := time.Now()
	s := NullTime{
		Valid: true,
		Time:  foo,
	}
	bits, err := s.MarshalJSON()
	checkError(err)
	fmt.Println(string(bits))
	expected, _ := foo.MarshalJSON()
	if !bytes.Equal(bits, expected) {
		t.Errorf("expected json marshal to be %s, was %s", string(bits), string(expected))
	}

	s = NullTime{
		Valid: false,
		Time:  foo,
	}
	bits, err = s.MarshalJSON()
	checkError(err)
	if !bytes.Equal(bits, []byte{}) {
		t.Errorf("expected json marshal to be %s, was %s", string(bits), string(expected))
	}
}
