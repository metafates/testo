package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"

	"testman"

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

type PluginWhichSkipsRandomTests struct{ *testman.T }

func (p PluginWhichSkipsRandomTests) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Hooks: plugin.Hooks{
			BeforeEach: func() {
				if rand.Int()%2 == 0 {
					p.Skip("random chose so")
				}
			},
		},
	}
}

type PluginWhichOverridesLog struct{ *testman.T }

func (p PluginWhichOverridesLog) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Overrides: plugin.Overrides{
			Log: func(plugin.FuncLog) plugin.FuncLog {
				return func(args ...any) {
					fmt.Printf("✨ %s\n", fmt.Sprint(args...))
				}
			},
			Logf: func(plugin.FuncLogf) plugin.FuncLogf {
				return func(format string, args ...any) {
					fmt.Printf("✨ %s\n", fmt.Sprintf(format, args...))
				}
			},
			Skip: func(plugin.FuncSkip) plugin.FuncSkip {
				return func(args ...any) {
					fmt.Printf("⚠️ Skipping because %s\n", fmt.Sprint(args...))

					p.SkipNow()
				}
			},
		},
	}
}
