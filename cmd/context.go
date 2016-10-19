package cmd

import (
	"fmt"

	accord "github.com/datascienceinc/accord/pkg"
	"github.com/datascienceinc/accord/pkg/module"
)

// Context is the struct that will contain all data for the process.
type Context struct {
	Tree *module.Tree
}

func NewContext(tr *module.Tree) *Context {
	return &Context{tr}
}

func (c *Context) ProcessEndpoints(f func(*accord.Endpoint)) {
	for name, child := range c.Tree.Children() {
		for _, ep := range child.Config().Endpoints {
			fmt.Printf("Module: %s -> %s\n", name, ep.URI)
		}
	}
}
