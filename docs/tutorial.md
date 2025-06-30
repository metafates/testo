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

`testo` uses its own `testing.T` wrapper to extend its capabilities.
All the methods of regular `testing.T` are available to use.
Moreover, this wrapper is compatible with `testing.TB` interface since it _embeds_ `testing.T`.

**Important**: Lets make an alias for it to avoid repetitions.

Create `main_test.go` file and put the following code in it:

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

To make `testo` compatible with `go test` invoke it from a regular testing function.

```go
func Test(t *testing.T) {
    // testo only needs to know what suite we have to run and what T does it use.
    testo.RunSuite[*Suite, T](t)
}
```

**Important:** specified `T` in `RunSuite` must match the `T`'s in tests arguments.
If not - error is raised _before_ running any tests.
Moreover, suite type must be a pointer.

## Running the tests

Run the tests as usual:

```go
go test . -v
```

You would see the following output:

```txt
=== RUN   Test
=== RUN   Test/Suite
=== RUN   Test/Suite/testo!
=== RUN   Test/Suite/testo!/TestAdd
--- PASS: Test (0.00s)
    --- PASS: Test/Suite (0.00s)
        --- PASS: Test/Suite/testo! (0.00s)
            --- PASS: Test/Suite/testo!/TestAdd (0.00s)
PASS
```

Notice the special `testo!` test - `testo` defines it internally to
make parallel tests work correctly with hooks.

You shouldn't worry about, as it doesn't affect your tests.
For example in `func (Suite) TestAdd(t T) { ... }` calling `t.Name()` would return `Test/Suite/TestAdd`.

## Suite hooks

Suite can define the following hooks:

- `BeforeAll(T)` - called before _all_ tests once. Passed `T` refers to the top-level test, for example `Test/Suite`.
- `BeforeEach(T)` - called before _each_ test. Passed `T` is the same as in actual test to run.
- `AfterEach(T)` - called after _each_ test is finished, but before cleanup. Passed `T` is the same as in actual test.
- `AfterAll(T)` - called after _all_ tests are finished once. It waits for all parallel tests to finish before running. Passed `T` refers to the top-level test, for example `Test/Suite`.

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

```txt
=== RUN   Test
=== RUN   Test/Suite
=== RUN   Test/Suite/testo!
=== RUN   Test/Suite/testo!/Add
    main_test.go:20: Starting: Test/Suite/TestAdd
    main_test.go:28: Test/Suite/TestAdd
    main_test.go:24: Finished: Test/Suite/TestAdd
--- PASS: Test (0.00s)
    --- PASS: Test/Suite (0.00s)
        --- PASS: Test/Suite/testo! (0.00s)
            --- PASS: Test/Suite/testo!/TestAdd (0.00s)
PASS
```

</details>

## Parametrized tests

You may want to test your code on various inputs to ensure it covers all cases.

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

We also have to define which parameters are passed to it.
To do so, define value providing methods.
They are named `CasesXXX` where `XXX` is the name of the parameter as specified by a field name.

```go
func (Suite) CasesA() []int {
    return []int{1, 2, 3, 4, 5}
}

func (Suite) CasesB() []int {
    return []int{11, 1000, 13}
}
```

Parametrized functions are called with the Cartesian product
of all values provided by `Cases` functions - 15 different cases in our example.

If test specifies a parameter for which `Cases` function does
not exist an error is raised before running any tests during static analysis.
The same goes for type mismatch.

## Sub-tests

Since tests in `testo` does not use `testing.T` directly,
running sub-tests (`t.Run`) is done differently to preserve `T` type inside sub-tests.

Helper function `testo.Run` is available for that.
It accepts `T` instance, sub-test name and an actual sub-test as a function.

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

For example:

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

Remember an alias for `T` defined previously:

```go
type T = *testo.T
```

Now add (install) plugins to it like that:

```go
type T = *struct{
    *testo.T

    ReverseTestsOrder
    OverrideLog
    AddNewMethods
    Timer
}
```

The only change needed. `testo` automatically initializes and uses specified plugins.

Since `AddNewMethods` is now embedded, it's possible to use new methods it defines:

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
