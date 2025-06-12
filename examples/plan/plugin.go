package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

type PluginWhichReversesTestOrder struct{}

// defining method "New" is optional - here we don't need it

func (PluginWhichReversesTestOrder) Plugin() plugin.Spec {
	return plugin.Spec{
		Plan: plugin.Plan{
			Modify: func(tests *[]plugin.PlannedTest) {
				slices.Reverse(*tests)
			},
		},
	}
}

type PluginWhichSkipsRandomTests struct{ *testo.T }

func (p PluginWhichSkipsRandomTests) Plugin() plugin.Spec {
	return plugin.Spec{
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

type PluginWhichOverridesLog struct{ *testo.T }

func (p PluginWhichOverridesLog) Plugin() plugin.Spec {
	return plugin.Spec{
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

type Timer struct {
	*testo.T

	start time.Time
}

func (t *Timer) Plugin() plugin.Spec {
	return plugin.Spec{
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
