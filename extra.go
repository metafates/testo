package testo

import "github.com/metafates/testo/plugin"

// Inspect returns meta information about given t.
//
// Note that all plugins and suite tests share
// the same pointer to the underlying [T].
func Inspect[T CommonT](t T) plugin.TInfo {
	return t.unwrap().info
}
