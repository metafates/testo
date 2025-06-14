// Package testifyplugin provides common wrappers for testify.
package testifyplugin

import (
	"github.com/metafates/testo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testify plugin providing utility helpers.
type Testify struct {
	*testo.T
}

// Require returns a new assertions object.
func (a *Testify) Require() *require.Assertions {
	return require.New(a)
}

// Assert returns a new assertions object.
func (a *Testify) Assert() *assert.Assertions {
	return assert.New(a)
}
