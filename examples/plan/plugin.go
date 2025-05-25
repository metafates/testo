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

type PluginWhichAddsNewTests struct{}

func (PluginWhichAddsNewTests) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Plan: plugin.Plan{
			Add: func() []plugin.Test {
				return []plugin.Test{
					{
						Name: "this test was added from plugin",
						Run: func(t plugin.T) {
							t.Log("Hello from virtual test!")
						},
					},
				}
			},
		},
	}
}
