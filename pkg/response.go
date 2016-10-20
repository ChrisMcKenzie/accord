package accord

import "net/http"

// Response ...
type Response struct {
	Headers http.Header
	Code    int         `hcl:"code"`
	Body    interface{} `hcl:"body"`
}

type Request struct {
	Headers http.Header
	Body    interface{} `hcl:"body"`
}
