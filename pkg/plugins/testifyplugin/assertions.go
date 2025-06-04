package testifyplugin

import (
	"github.com/metafates/tego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Testify struct {
	*tego.T
}

func (a *Testify) Require() *require.Assertions {
	return require.New(a)
}

func (a *Testify) Assert() *assert.Assertions {
	return assert.New(a)
}
