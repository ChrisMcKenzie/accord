package module

import getter "github.com/hashicorp/go-getter"

func getStorage(s getter.Storage, key string, src string) (string, bool, error) {
	// Get the directory where the module is.
	return s.Dir(key)
}
