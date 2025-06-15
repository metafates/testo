//go:build example

package main

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

type T = *struct {
	*testo.T

	ReverseTestsOrder
	OverrideLog
	AddNewMethods
	Timer
}

func Test(t *testing.T) {
	// testo only needs to know what suite we have to run and what T does it use.
	testo.RunSuite[*Suite, T](t)
}

type Suite struct{}

func (Suite) BeforeEach(t T) {
	t.Logf("Starting: %s", t.Name())
}

func (Suite) AfterEach(t T) {
	t.Logf("Finished: %s", t.Name())
}

// Suite tests are just a methods and they follow the same naming rules, as regular tests are.
// That is, they must have "Test" prefix. Also, they must use our custom T we defined earlier.
func (Suite) TestAdd(t T) {
	// remember - T has all the methods regular testing.T does.
	if Add(2, 2) != 4 {
		t.Fatal("2 + 2 must equal 4")
	}
}

func (Suite) CasesA() []int {
	return []int{1, 2, 3, 4, 5}
}

func (Suite) CasesB() []int {
	return []int{11, 1000, 13}
}

func (Suite) CasesC() []int {
	return []int{-4, -99, 9}
}

func (Suite) TestAddButParametrized(t T, params struct{ A, B, C int }) {
	testo.Run(t, "commutative", func(t T) {
		if Add(params.A, params.B) != Add(params.B, params.A) {
			t.Errorf("%[1]d + %[2]d != %[2]d + %[1]d", params.A, params.B)
		}
	})

	testo.Run(t, "associative", func(t T) {
		if Add(Add(params.A, params.B), params.C) != Add(params.A, Add(params.B, params.C)) {
			t.Errorf(
				"(%[1]d + %[2]d) + %[3]d != %[1]d + (%[2]d + %[3]d)",
				params.A,
				params.B,
				params.C,
			)
		}
	})
}

type ReverseTestsOrder struct{}

// plugins can implement this function to provide
// certain plugin functionality.
//
// It is optional - see AddNewMethods plugin.
func (ReverseTestsOrder) Plugin() plugin.Spec {
	return plugin.Spec{
		Plan: plugin.Plan{
			Modify: func(tests *[]plugin.PlannedTest) {
				slices.Reverse(*tests)
			},
		},
	}
}

type OverrideLog struct{}

func (OverrideLog) Plugin() plugin.Spec {
	return plugin.Spec{
		Overrides: plugin.Overrides{
			Log: func(f plugin.FuncLog) plugin.FuncLog {
				return func(args ...any) {
					// this will be printed each time t.Log is called.
					fmt.Println("Inside log override")
					f(args...)
				}
			},
		},
	}
}

// we can embed testo.T in plugins - it will be automatically initialized
// and share the same testo.T as an actual T from the current test.
type AddNewMethods struct{ *testo.T }

// you will see later how we can access this function in tests.
func (a AddNewMethods) Explode() { a.Fatal("BOOM") }

type Timer struct {
	*testo.T
	start time.Time
}

func (t *Timer) Plugin() plugin.Spec {
	return plugin.Spec{
		Hooks: plugin.Hooks{
			BeforeEach: plugin.Hook{
				Priority: plugin.TryLast,
				Func: func() {
					// .Plugin() is called for each test, therefore
					// we can modify Timer fields safely (new instance for each test).
					t.start = time.Now()
				},
			},
			AfterEach: plugin.Hook{
				Priority: plugin.TryFirst,
				Func: func() {
					elapsed := time.Since(t.start)

					fmt.Printf("Test %q took %s\n", t.Name(), elapsed)
				},
			},
		},
	}
}
