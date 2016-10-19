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
	c.processEndpoints(c.Tree, f)

	for _, child := range c.Tree.Children() {
		c.processEndpoints(child, f)
	}
}

func (c *Context) processEndpoints(tr *module.Tree, f func(*accord.Endpoint)) {
	for _, ep := range tr.Config().Endpoints {
		fmt.Printf("Module %s:\n", tr.Name())
		f(ep)
	}
}
