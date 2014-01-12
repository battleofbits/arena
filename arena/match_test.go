package arena

import (
	"bytes"
	"encoding/json"
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
}

func TestNullTimeNullJSON(t *testing.T) {
	foo := time.Now()
	s := NullTime{
		Valid: false,
		Time:  foo,
	}
	bits, err := s.MarshalJSON()
	checkError(err)
	if !bytes.Equal(bits, []byte{}) {
		t.Errorf("expected json marshal to be empty, was %s", string(bits))
	}
}

func TestNullStringFullRound(t *testing.T) {
	s := NullString{
		Valid:  true,
		String: "foo",
	}
	bits, err := json.Marshal(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	var output NullString
	err = json.Unmarshal(bits, &output)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(output)
	if output.String != "foo" {
		t.Errorf("Expected NullString's String to be foo, was %s", output.String)
	}
}
