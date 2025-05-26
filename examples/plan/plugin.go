package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"time"

	"testman"

	"testman/plugin"
)

type PluginWhichReversesTestOrder struct{}

// defining method "New" is optional - here we don't need it

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
							// this log can be be overridden by other plugins
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
			BeforeEach: plugin.Hook{
				Func: func() {
					if rand.Int()%2 == 0 {
						p.Skip("random chose so")
					}
				},
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

type PluginWhichRenamesTests struct{}

func (PluginWhichRenamesTests) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Plan: plugin.Plan{
			Rename: func(name string) string {
				return name + " [renamed]"
			},
		},
	}
}

type Timer struct {
	*testman.T

	start time.Time
}

func (t *Timer) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Hooks: plugin.Hooks{
			BeforeEach: plugin.Hook{
				Priority: plugin.TryLast, // instruct to run this hook as late as possible
				Func: func() {
					t.start = time.Now()
				},
			},
			AfterEach: plugin.Hook{
				Priority: plugin.TryFirst, // and this hook to be run as early as possible
				Func: func() {
					duration := time.Since(t.start)

					fmt.Printf("⌛ Test %q took %s to complete\n", t.Name(), duration)
				},
			},
		},
	}
}
