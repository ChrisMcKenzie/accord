package accord

// Config defines the hcl config of an accord file.
type Config struct {
	Dir       string
	Endpoints []*Endpoint
	Modules   []*Module
}

// Module represents an import of a remote accord
type Module struct {
	Name   string
	Source string
}
