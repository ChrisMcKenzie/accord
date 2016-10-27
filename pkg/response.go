package accord

import "net/http"

// Response ...
type Response struct {
	Headers http.Header
	Code    int         `hcl:"code"`
	Body    interface{} `hcl:"body"`
}

// Request ...
type Request struct {
	Headers http.Header
	Body    interface{}       `hcl:"body"`
	Query   map[string]string `hcl:"query"`
}
