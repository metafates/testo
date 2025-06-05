package main

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ShuffleTests struct{}

func (ShuffleTests) Plugin() plugin.Spec {
	return plugin.Spec{
		Plan: plugin.Plan{
			Modify: func(tests []plugin.PlannedTest) []plugin.PlannedTest {
				slices.SortFunc(
					tests,
					func(_, _ plugin.PlannedTest) int { return rand.IntN(3) - 1 },
				)

				// modify receives a slice clone, so modifying it in-place is not enough.
				// we must return a new slice
				return tests
			},
		},
	}
}

type MakeLogsPretty struct{}

func (MakeLogsPretty) Plugin() plugin.Spec {
	return plugin.Spec{
		Overrides: plugin.Overrides{
			Log: func(f plugin.FuncLog) plugin.FuncLog {
				return func(args ...any) {
					f("✨ " + fmt.Sprint(args...))
				}
			},
			Logf: func(f plugin.FuncLogf) plugin.FuncLogf {
				return func(format string, args ...any) {
					f("✨ "+format, args...)
				}
			},
		},
	}
}

// *testo.T inside plugins will initialized automatically
// and will refer to the current test
type AnnounceTests struct{ *testo.T }

// this is invoked for at each hook stage, so *testo.T will always refer to the current test.
func (a AnnounceTests) Plugin() plugin.Spec {
	return plugin.Spec{
		Hooks: plugin.Hooks{
			BeforeEach: plugin.Hook{
				Priority: plugin.TryFirst, // we can set order priority
				Func: func() {
					a.Logf("Test %q started", a.Name())
				},
			},
			AfterEach: plugin.Hook{
				Priority: plugin.TryLast, // priority is just an int, so we can write 9999 instead
				Func: func() {
					a.Logf("Test %q finished", a.Name())
				},
			},
		},
	}
}

// plugins are not required to implement Plugin() method.
// this plugin just provides some helper methods.
type Testify struct{ *testo.T }

func (t Testify) Require() *require.Assertions { return require.New(t) }
func (t Testify) Assert() *assert.Assertions   { return assert.New(t) }
