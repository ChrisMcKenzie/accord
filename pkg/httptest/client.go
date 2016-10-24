package httptest

import (
	"net/http"

	accord "github.com/datascienceinc/accord/pkg"
)

// Client that tests that the expected response is given
type Client struct {
	*http.Client
}

// NewClient create a new instance of *httptest.Client
func NewClient() *Client {
	return &Client{&http.Client{}}
}

// Evaluate will call http.Client.Do and evaluate the response
func (c *Client) Evaluate(req *http.Request, expected *accord.Response) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	// first things first is the status code right?
	if resp.StatusCode != expected.Code {
		return diffError(string(expected.Code), string(resp.StatusCode))
	}

	// Lets go through all the headers in the expectation and make sure they match
	// we dont want to test that all headers match and exist because some servers
	// will send extra headers (ie. X-Powered-By)
	for h := range expected.Headers {
		if expected.Headers.Get(h) != resp.Header.Get(h) {
			return diffError(expected.Headers.Get(h), resp.Header.Get(h))
		}
	}

	// Lets check if the bodies are equal now!
	err = compareResponse(resp, expected)

	return err
}
