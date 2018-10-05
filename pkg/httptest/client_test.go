package httptest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	accord "github.com/ChrisMcKenzie/accord/pkg"
)

func httpHandler(method string, res *accord.Response, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			t.Errorf("Expected method to be %s, got %s", method, req.Method)
		}

		for key := range res.Headers {
			w.Header().Set(key, res.Headers.Get(key))
		}

		w.WriteHeader(res.Code)

		resp := parseBody(res.Headers, res.Body)
		w.Write(resp.Bytes())
	}
}

func newRequest(method, url string, req *accord.Request) (*http.Request, error) {
	var buf bytes.Buffer
	if req != nil {
		buf = parseBody(req.Headers, req.Body)
	}

	r, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func TestHttpClient(t *testing.T) {
	cases := []struct {
		name     string
		errNil   bool
		method   string
		url      string
		request  *accord.Request
		response *accord.Response
		expected *accord.Response
	}{
		{
			name:     "Basic Request",
			errNil:   true,
			method:   "GET",
			url:      "/test",
			request:  &accord.Request{},
			response: &accord.Response{},
			expected: &accord.Response{},
		},
		{
			name:     "Failing Status Code request",
			errNil:   false,
			method:   "GET",
			url:      "/test",
			request:  &accord.Request{},
			response: &accord.Response{},
			expected: &accord.Response{
				Code: 400,
			},
		},
		{
			name:     "Failing Header request",
			errNil:   false,
			method:   "GET",
			url:      "/test",
			request:  &accord.Request{},
			response: &accord.Response{},
			expected: &accord.Response{
				Headers: http.Header(map[string][]string{
					"X-My-Header": []string{"test"},
				}),
			},
		},
		{
			name:     "Failing Body request",
			errNil:   false,
			method:   "GET",
			url:      "/test",
			request:  &accord.Request{},
			response: &accord.Response{},
			expected: &accord.Response{
				Body: "test",
			},
		},
		{
			name:    "JSON body request",
			errNil:  true,
			method:  "GET",
			url:     "/test",
			request: &accord.Request{},
			response: &accord.Response{
				Headers: http.Header(map[string][]string{
					"Content-Type": []string{"application/json"},
				}),
				Body: map[string]string{
					"test": "yes",
				},
			},
			expected: &accord.Response{
				Headers: http.Header(map[string][]string{
					"Content-Type": []string{"application/json"},
				}),
				Body: map[string]string{
					"test": "yes",
				},
			},
		},
	}

	client := NewClient()
	for _, c := range cases {
		fmt.Printf("==> %s\n", c.name)
		server := httptest.NewServer(
			http.HandlerFunc(
				httpHandler(
					c.method,
					c.response,
					t,
				),
			),
		)
		defer server.Close()

		req, err := newRequest(
			c.method,
			strings.Join([]string{server.URL, c.url}, "/"),
			c.request)
		if err != nil {
			t.Error(err)
		}

		err = client.Evaluate(req, c.expected)
		if c.errNil && err != nil {
			t.Errorf("Expected evaluate to return nil got:\n%s\n", err)
		}
	}
}
