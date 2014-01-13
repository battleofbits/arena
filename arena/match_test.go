package arena

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestNullTimeJSON(t *testing.T) {
	t.Parallel()
	foo := time.Now()
	s := NullTime{
		Valid: true,
		Time:  foo,
	}
	bits, err := json.Marshal(s)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected, _ := json.Marshal(foo)
	if !bytes.Equal(bits, expected) {
		t.Errorf("expected json marshal to be %s, was %s", string(bits),
			string(expected))
	}
}

func TestNullTimeNullJSON(t *testing.T) {
	t.Parallel()
	s := NullTime{
		Valid: false,
	}
	bits, err := json.Marshal(s)
	checkError(err)
	if !bytes.Equal(bits, []byte("null")) {
		t.Errorf("expected json marshal to be empty, was %s", string(bits))
	}
}

func TestNullStringFullRound(t *testing.T) {
	t.Parallel()
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
	if output.String != "foo" {
		t.Errorf("Expected NullString's String to be foo, was %s", output.String)
	}
}
