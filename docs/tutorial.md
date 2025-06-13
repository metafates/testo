# Tutorial

Take a guided tour of `testo` by making a simple plugins and running the tests using various features.

## Create a new Go project

```bash
mkdir testo-tutorial
cd testo-tutorial
go mod init testo-tutorial
```

Create `main.go` and put the following code in it:

```go
package main

// Add returns sum of the given integers.
func Add(a, b int) int { return a + b }

func main() {}
```

## Writing tests

`testo` uses its own `testing.T` wrapper to extend its functionality.
All the methods of regular `testing.T` are available to use.
Moreover, this wrapper is compatible with `testing.TB` interface since it embeds `testing.T`.

Lets make an alias for it to avoid repetitions (this is important, we'll modify it later).

Create `main_test.go` file and put the following code in it.

```go
package main

import (
	"testing"

    "github.com/metafates/testo"
)

type T = *testo.T
```

We also need a suite. A suite is anything that can define its own methods.

Let's start with an empty struct.

```go
type Suite struct{}

// Suite tests are just a methods and they follow the same naming rules, as regular tests are.
// That is, they must have "Test" prefix. Also, they must use our custom T we defined earlier.
func (Suite) TestAdd(t T) {
    // remember - T has all the methods regular testing.T does.
    if Add(2, 2) != 4 {
        t.Fatal("2 + 2 must equal 4")
    }
}
```

To make `testo` compatible with `go test` we invoke it from a regular testing function.

```go
func Test(t *testing.T) {
    // testo only needs to know what suite we have to run and what T does it use.
    testo.RunSuite[Suite, T](t)
}
```

**Important:** specified `T` in `RunSuite` must match the `T`'s in tests arguments.
If not - error is raised *before* running any tests.

## Running the tests

Nothing new here:

```go
go test . -v
```

You would see the following output:

```
=== RUN   Test
=== RUN   Test/Suite
=== RUN   Test/Suite/testo!
=== RUN   Test/Suite/testo!/Add
--- PASS: Test (0.00s)
    --- PASS: Test/Suite (0.00s)
        --- PASS: Test/Suite/testo! (0.00s)
            --- PASS: Test/Suite/testo!/Add (0.00s)
PASS
```

Notice the special `testo!` test - `testo` defines it internally to
make parallel tests work correctly with hooks.

You should not worry about, as it does not affect your tests.
For example, `t.Name()` method would remove it for you, e.g. in
`func (Suite) TestAdd(t T) { ... }` calling `t.Name()` would return `Test/Suite/Add`.

## Suite hooks

Suite can define the following hooks:

- `BeforeAll(T)` - called before *all* tests once. Passed `T` refers to the top-level test (e.g. `Test/Suite`).
- `BeforeEach(T)` - called before *each* test. Passed `T` is the same as in actual test to be run.
- `AfterEach(T)` - called after *each* test is finished (but before cleanup). Passed `T` is the same as in actual test.
- `AfterAll(T)` - called after *all* tests are finished once. It will wait for all parallel tests to finish before running. Passed `T` refers to the top-level test (e.g. `Test/Suite`)

Hooks are defined as suite methods:

```go
func (Suite) BeforeEach(t T) {
    t.Logf("Starting: %s", t.Name())
}

func (Suite) AfterEach(t T) {
    t.Logf("Finished: %s", t.Name())
}
```

<details>
<summary>Output:</summary>

```
=== RUN   Test
=== RUN   Test/Suite
=== RUN   Test/Suite/testo!
=== RUN   Test/Suite/testo!/Add
    main_test.go:20: Starting: Test/Suite/Add
    main_test.go:28: Test/Suite/Add
    main_test.go:24: Finished: Test/Suite/Add
--- PASS: Test (0.00s)
    --- PASS: Test/Suite (0.00s)
        --- PASS: Test/Suite/testo! (0.00s)
            --- PASS: Test/Suite/testo!/Add (0.00s)
PASS
```

</details>

## Parametrized tests

You may want to test your code on various inputs to ensure all cases are covered.

`testo` makes it easier by letting you define parametrized tests.

Parametrized tests follow the same naming rules as regular tests, but accept a second argument after `T`
as a struct of required parameters. This struct can be anonymous or named.

```go
func (Suite) TestAddButParametrized(t T, params struct{ A, B int }) {
    if Add(params.A, params.B) != Add(params.B, params.A) {
        t.Errorf("%[1]d + %[2]d != %[2]d + %[1]d", params.A, params.B)
    }
}
```

But we also have to define which parameters it will use.
To do so, we need to define value providing methods.
They are named `CasesXXX` where `XXX` is the name of the parameter as specified by a field name.

```go
func (Suite) CasesA() []int {
    return []int{1, 2, 3, 4, 5}
}

func (Suite) CasesB() []int {
    return []int{11, 1000, 13}
}
```

Parametrized functions will run with the Cartesian product
of all values provided by `Cases` functions - 15 different cases in our example.

If test will specify a parameter for which `Cases` function does
not exist an error will be raised before running any tests during static analysis.
The same goes for type mismatch.

## Subtests

Since tests in `testo` does not use `testing.T` directly,
running subtests (`t.Run`) is done differently to preserve `T` type inside subtests.

Helper function `testo.Run` is available for that.
It accepts `T` instance, subtest name and an actual subtest as a function.

```go
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
            t.Errorf("(%[1]d + %[2]d) + %[3]d != %[1]d + (%[2]d + %[3]d)", params.A, params.B, params.C)
        }
    })
}
```

## Plugins

One of the biggest features of `testo` is plugin system.

### Writing plugins

Plugins are anything that can define its own methods.

Let's make some simple plugins

1. A plugin which reverses the order of tests.
2. A plugin which overrides `t.Log` function.
3. A plugin which adds new methods to `T`.
4. A plugin which shows time taken for each test.

```go
import (
	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

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

type OverrideLog struct {}

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
```

### Using plugins

We defined an alias for `T` previously:

```go
type T = *testo.T
```

Now we can add (install) plugins to it like that:

```go
type T = *struct{
    *testo.T

    ReverseTestsOrder
    OverrideLog
    AddNewMethods
    Timer
}
```

This is the only change needed. `testo` will automatically initialize and use them.

Since `AddNewMethods` is now embedded, we can use new methods it defines:

```go
func (Suite) TestBoom(t T) {
    t.Explode()
}
```

## Putting it all together

```go
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
	testo.RunSuite[Suite, T](t)
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
			t.Errorf("(%[1]d + %[2]d) + %[3]d != %[1]d + (%[2]d + %[3]d)", params.A, params.B, params.C)
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
```
