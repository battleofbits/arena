package arena

import (
	"bytes"
	"net/http"
)

// Make a request to a client's server
func MakeRequest(url string, postBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(postBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "battleofbits/0.1")
	// XXX, set a timeout here
	return client.Do(req)
}
