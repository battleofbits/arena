package arena

import (
	"bytes"
	"fmt"
	"net/http"
)

const VERSION = "0.1"

// Make a request to a client's server
func MakeRequest(url string, postBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(postBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", fmt.Sprintf("battleofbits/%s", VERSION))
	// XXX, set a timeout here
	return client.Do(req)
}
