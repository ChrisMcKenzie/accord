package module

// Code Borrowed from terraform mostly

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	accord "github.com/ChrisMcKenzie/accord/pkg"
	getter "github.com/hashicorp/go-getter"
)

const rootName = "root"

// Tree is a struct containing a root config and all of the child modules.
type Tree struct {
	name     string
	config   *accord.Config
	children map[string]*Tree
	path     []string
	lock     sync.RWMutex
}

// NewTree returns a new Tree for the given config structure.
func NewTree(name string, c *accord.Config) *Tree {
	return &Tree{config: c, name: name}
}

// NewEmptyTree returns a new tree that is empty (contains no configuration).
func NewEmptyTree() *Tree {
	t := &Tree{config: &accord.Config{}}

	// We do this dummy load so that the tree is marked as "loaded". It
	// should never fail because this is just about a no-op. If it does fail
	// we panic so we can know its a bug.
	if err := t.Load(nil); err != nil {
		panic(err)
	}

	return t
}

// NewTreeModule is like NewTree except it parses the configuration in
// the directory and gives it a specific name. Use a blank name "" to specify
// the root module.
func NewTreeModule(name, path string) (*Tree, error) {
	c, err := accord.LoadConfig(path)
	if err != nil {
		return nil, err
	}

	return NewTree(name, c), nil
}

// Config ...
func (t *Tree) Config() *accord.Config {
	return t.config
}

// Child returns the child with the given path (by name).
func (t *Tree) Child(path []string) *Tree {
	if len(path) == 0 {
		return t
	}

	c := t.Children()[path[0]]
	if c == nil {
		return nil
	}

	return c.Child(path[1:])
}

// Children returns the children of this tree (the modules that are
// imported by this root).
//
// This will only return a non-nil value after Load is called.
func (t *Tree) Children() map[string]*Tree {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.children
}

// Name returns the name of the tree. This will be "<root>" for the root
// tree and then the module name given for any children.
func (t *Tree) Name() string {
	if t.name == "" {
		return rootName
	}

	return t.name
}

// Modules returns the list of modules that this tree imports.
//
// This is only the imports of _this_ level of the tree. To retrieve the
// full nested imports, you'll have to traverse the tree.
func (t *Tree) Modules() []*accord.Module {
	return t.config.Modules
}

// Load loads the configuration of the entire tree.
//
// The parameters are used to tell the tree where to find modules and
// whether it can download/update modules along the way.
//
// Calling this multiple times will reload the tree.
//
// Various semantic-like checks are made along the way of loading since
// module trees inherently require the configuration to be in a reasonably
// sane state: no circular dependencies, proper module sources, etc. A full
// suite of validations can be done by running Validate (after loading).
func (t *Tree) Load(s getter.Storage) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	// Reset the children if we have any
	t.children = nil

	modules := t.Modules()
	children := make(map[string]*Tree)

	// Go through all the modules and get the directory for them.
	for _, m := range modules {
		if _, ok := children[m.Name]; ok {
			return fmt.Errorf(
				"module %s: duplicated. module names must be unique", m.Name)
		}

		// Determine the path to this child
		path := make([]string, len(t.path), len(t.path)+1)
		copy(path, t.path)
		path = append(path, m.Name)

		// Split out the subdir if we have one
		source, subDir := getter.SourceDirSubdir(m.Source)

		source, err := getter.Detect(source, t.config.Dir, getter.Detectors)
		if err != nil {
			return fmt.Errorf("module %s: %s", m.Name, err)
		}

		// Check if the detector introduced something new.
		source, subDir2 := getter.SourceDirSubdir(source)
		if subDir2 != "" {
			subDir = filepath.Join(subDir2, subDir)
		}

		// Get the directory where this module is so we can load it
		key := strings.Join(path, ".")
		key = fmt.Sprintf("root.%s-%s", key, m.Source)
		dir, ok, err := getStorage(s, key, source)
		if err != nil {
			return err
		}
		if !ok {
			s.Get(key, source, true)
		}

		// If we have a subdirectory, then merge that in
		if subDir != "" {
			dir = filepath.Join(dir, subDir)
		}

		// Load the configurations.Dir(source)
		children[m.Name], err = NewTreeModule(m.Name, dir)
		if err != nil {
			return fmt.Errorf(
				"module %s: %s", m.Name, err)
		}

		// Set the path of this child
		children[m.Name].path = path
	}

	// Go through all the children and load them.
	for _, c := range children {
		if err := c.Load(s); err != nil {
			return err
		}
	}

	// Set our tree up
	t.children = children

	return nil
}
