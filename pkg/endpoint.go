package accord

// Endpoint Defines the data for an endpoint in an accord
type Endpoint struct {
	URI      string    `hcl:",key"`
	Method   string    `hcl:"method"`
	Response *Response `hcl:"response"`
	Request  *Request  `hcl:"request"`
}
