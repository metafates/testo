# Testo

<img src="https://github.com/user-attachments/assets/66844de4-4b13-428a-b924-1f26718cee41" align="right" width="250" alt="testo mascot">

Testo is a modular testing framework built on top of `testing.T`.

## Features

- Test suites
- Plugins with hooks, test planning, built-in function overrides (Log, Error)
- Parametrized tests
- Parallel tests

## Quick Start

```go
package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type T struct {
	*testo.T

	// This is how plugins are "installed"
	ShuffleTests
	MakeLogsPretty
	AnnounceTests
	Testify
}

// We need to write a regular test function to bridge
// testo with "go test"
func Test(t *testing.T) {
	testo.RunSuite[*Suite, *T](t)
}

// Then we define our suite. It can be struct and include some fields or any other
// type as long as it can have its own methods.
type Suite struct {
	names []string
}

// ================= HOOKS =========================

// This hook will run before all tests.
func (s *Suite) BeforeAll(t *T) {
	// this function simulates some expensive operation,
	// like fetching from DB.
	s.names = getNames()
}

// This hook will run before each test.
// It is run in the same context as the test itself.
// That is, it has the same pointer to T as the actual test.
//
// Using that, we can make all tests parallel just by running t.Parallel() here
func (s *Suite) BeforeEach(t *T) {}

// Same as BeforeEach but runs after each test.
// Again - the same pointer to T as the actual T.
func (s *Suite) AfterEach(t *T) {}

// This hook will run after all tests are finished.
// It will wait for all parallel tests to finish before running
func (s *Suite) AfterAll(t *T) {}

// ================ WRITING THE ACTUAL TESTS ==========================

// A basic test as you would write without any frameworks.
// Tests are methods with "Test..." prefix.
//
// NOTE: it uses the same type of T as we passed to RunSuite. In our case it's *T.
func (Suite) TestAdd(t *T) {
	if 2+2 != 4 {
		t.Fatal("2 + 2 must be equal 4")
	}
}

// we can functions from testing libraries like testify,
// because our T is compatible with common testing.T interface
func (s Suite) TestWithTestify(t *T) {
	require.Equal(t, 4, 2+2, "2 + 2 must be equal 4")

	// since we've embedded Testify plugin, we can write it like that
	t.Require().Equal(4, 2+2, "2 + 2 must be equal 4")
}

// we can also run subtests
func (s Suite) TestWithSubtests(t *T) {
	// however, we can't run t.Run because we won't be able to preserve our custom T type.
	// just use testo.Run helper to run subtests.
	//
	// BeforeEach, AfterEach won't be triggered for subtests.
	testo.Run(t, "subtest", func(t *T) {
		t.Log("We are inside subtest now!")

		testo.Run(t, "nested", func(t *T) {
			t.Log("Subtest in a subtest")
		})
	})
}

// That's how we can define parametrized tests.
// Their signature is the same as regular tests but it has one
// additional argument of struct type (anonymous or named).
// Public fields of this struct define which parameters this test accept.
// Parameter values come from "CasesX" methods.
// This function will be called with all combinations of param values.
func (Suite) TestWithSomeParams(t *T, params struct{ Foo, Bar string }) {
	t.Logf("Foo is %q and bar is %q", params.Foo, params.Bar)
}

func (s Suite) CasesFoo() []string {
	// it was set in BeforeAll hook
	return s.names
}

func (Suite) CasesBar() []string {
	return []string{"abc", "xyz"}
}

// Parallel tests are supported.
// Since BeforeEach and AfterEach run in the same test context,
// They will be properly executed.
func (Suite) TestWithParallel(t *T) {
	t.Parallel() // this is okay and will make this test parallel as expected.

	// One limitation though
	testo.Run(t, "nested parallel", func(t *T) {
		// Running parallel at this level (one level nested) is
		// not supported and will be ignored.
		// This is technical limitation of "go test" works, which would affect hooks order.
		t.Parallel()

		testo.Run(t, "nested-nested parallel", func(t *T) {
			// but parallel tests at this level are OK and SUPPORTED.
			t.Parallel()
		})

		testo.Run(t, "another nested-nested parallel", func(t *T) {
			t.Parallel()
		})
	})
}

// this function simulates some expensive operation like fetching from DB.
func getNames() []string {
	time.Sleep(2 * time.Second)

	return []string{"Alice", "John", "Bob"}
}

// ================= PLUGINS ==========================

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
```

<details>
<summary>Output log</summary>

```
=== RUN   Test
=== RUN   Test/Suite
=== RUN   Test/Suite/testo!
=== RUN   Test/Suite/testo!/WithParallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel" started
=== PAUSE Test/Suite/testo!/WithParallel
=== RUN   Test/Suite/testo!/Add
    main_test.go:181: ✨ Test "Test/Suite/Add" started
    main_test.go:181: ✨ Test "Test/Suite/Add" finished
=== RUN   Test/Suite/testo!/WithTestify
    main_test.go:181: ✨ Test "Test/Suite/WithTestify" started
    main_test.go:181: ✨ Test "Test/Suite/WithTestify" finished
=== RUN   Test/Suite/testo!/WithSomeParams_case_1
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_1" started
    main_test.go:181: ✨ Foo is "Alice" and bar is "abc"
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_1" finished
=== RUN   Test/Suite/testo!/WithSubtests
    main_test.go:181: ✨ Test "Test/Suite/WithSubtests" started
=== RUN   Test/Suite/testo!/WithSubtests/subtest
    main_test.go:181: ✨ Test "Test/Suite/WithSubtests/subtest" started
    main_test.go:176: ✨ We are inside subtest now!
=== RUN   Test/Suite/testo!/WithSubtests/subtest/nested
    main_test.go:181: ✨ Test "Test/Suite/WithSubtests/subtest/nested" started
    main_test.go:176: ✨ Subtest in a subtest
    main_test.go:181: ✨ Test "Test/Suite/WithSubtests/subtest/nested" finished
=== NAME  Test/Suite/testo!/WithSubtests/subtest
    main_test.go:181: ✨ Test "Test/Suite/WithSubtests/subtest" finished
=== NAME  Test/Suite/testo!/WithSubtests
    main_test.go:181: ✨ Test "Test/Suite/WithSubtests" finished
=== RUN   Test/Suite/testo!/WithSomeParams_case_2
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_2" started
    main_test.go:181: ✨ Foo is "John" and bar is "abc"
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_2" finished
=== RUN   Test/Suite/testo!/WithSomeParams_case_6
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_6" started
    main_test.go:181: ✨ Foo is "Bob" and bar is "xyz"
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_6" finished
=== RUN   Test/Suite/testo!/WithSomeParams_case_4
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_4" started
    main_test.go:181: ✨ Foo is "Alice" and bar is "xyz"
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_4" finished
=== RUN   Test/Suite/testo!/WithSomeParams_case_5
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_5" started
    main_test.go:181: ✨ Foo is "John" and bar is "xyz"
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_5" finished
=== RUN   Test/Suite/testo!/WithSomeParams_case_3
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_3" started
    main_test.go:181: ✨ Foo is "Bob" and bar is "abc"
    main_test.go:181: ✨ Test "Test/Suite/WithSomeParams_case_3" finished
=== CONT  Test/Suite/testo!/WithParallel
=== RUN   Test/Suite/testo!/WithParallel/nested_parallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel/nested_parallel" started
    main_test.go:176: ✨ WARN: running Parallel() at this level is not supported and will be ignored
=== RUN   Test/Suite/testo!/WithParallel/nested_parallel/nested-nested_parallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel/nested_parallel/nested-nested_parallel" started
=== PAUSE Test/Suite/testo!/WithParallel/nested_parallel/nested-nested_parallel
=== RUN   Test/Suite/testo!/WithParallel/nested_parallel/another_nested-nested_parallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel/nested_parallel/another_nested-nested_parallel" started
=== PAUSE Test/Suite/testo!/WithParallel/nested_parallel/another_nested-nested_parallel
=== NAME  Test/Suite/testo!/WithParallel/nested_parallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel/nested_parallel" finished
=== CONT  Test/Suite/testo!/WithParallel/nested_parallel/nested-nested_parallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel/nested_parallel/nested-nested_parallel" finished
=== CONT  Test/Suite/testo!/WithParallel/nested_parallel/another_nested-nested_parallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel/nested_parallel/another_nested-nested_parallel" finished
=== NAME  Test/Suite/testo!/WithParallel
    main_test.go:181: ✨ Test "Test/Suite/WithParallel" finished
--- PASS: Test (2.01s)
    --- PASS: Test/Suite (2.01s)
        --- PASS: Test/Suite/testo! (0.00s)
            --- PASS: Test/Suite/testo!/Add (0.00s)
            --- PASS: Test/Suite/testo!/WithTestify (0.00s)
            --- PASS: Test/Suite/testo!/WithSomeParams_case_1 (0.00s)
            --- PASS: Test/Suite/testo!/WithSubtests (0.00s)
                --- PASS: Test/Suite/testo!/WithSubtests/subtest (0.00s)
                    --- PASS: Test/Suite/testo!/WithSubtests/subtest/nested (0.00s)
            --- PASS: Test/Suite/testo!/WithSomeParams_case_2 (0.00s)
            --- PASS: Test/Suite/testo!/WithSomeParams_case_6 (0.00s)
            --- PASS: Test/Suite/testo!/WithSomeParams_case_4 (0.00s)
            --- PASS: Test/Suite/testo!/WithSomeParams_case_5 (0.00s)
            --- PASS: Test/Suite/testo!/WithSomeParams_case_3 (0.00s)
            --- PASS: Test/Suite/testo!/WithParallel (0.00s)
                --- PASS: Test/Suite/testo!/WithParallel/nested_parallel (0.00s)
                    --- PASS: Test/Suite/testo!/WithParallel/nested_parallel/nested-nested_parallel (0.00s)
                    --- PASS: Test/Suite/testo!/WithParallel/nested_parallel/another_nested-nested_parallel (0.00s)
PASS
ok      github.com/metafates/testo/examples/tour        2.186s
```
</details>
