package accord

// Response ...
type Response struct {
	Headers map[string]string `hcl:"headers"`
	Code    int               `hcl:"code"`
	Body    string            `hcl:"body"`
}

type Request struct {
	Header map[string]string `hcl:"headers"`
	Body   string            `hcl:"body"`
}
