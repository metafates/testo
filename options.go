package testo

import (
	"slices"
	"sync"

	"github.com/metafates/testo/plugin"
)

//nolint:gochecknoglobals // global variables are required in this case.
var (
	defaultOptions      []plugin.Option
	defaultOptionsMutex sync.RWMutex
)

// AddDefaultOptions adds given options to the global options.
//
// Global options are prepended to each [RunSuite] call.
func AddDefaultOptions(options ...plugin.Option) {
	defaultOptionsMutex.Lock()
	defer defaultOptionsMutex.Unlock()

	defaultOptions = append(defaultOptions, options...)
}

func getDefaultOptions() []plugin.Option {
	defaultOptionsMutex.RLock()
	defer defaultOptionsMutex.RUnlock()

	return slices.Clone(defaultOptions)
}
