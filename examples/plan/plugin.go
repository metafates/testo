package main

import (
	"cmp"

	"testman/plugin"
)

type PluginWhichReversesTestOrder struct{}

// defining New method is optional - here we don't need it

func (PluginWhichReversesTestOrder) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Plan: plugin.Plan{
			Sort: func(a, b string) int {
				return cmp.Compare(b, a)
			},
		},
	}
}
