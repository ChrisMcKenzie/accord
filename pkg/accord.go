package accord

// Accord defines the hcl config of an accord file.
type Accord struct {
	Endpoints []*Endpoint `hcl:"endpoint,expand"`
	Source    string      `hcl:"source"`
}
