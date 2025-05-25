# Testman

A Go testing framework with support for **test suites**, **plugins**, and **test hooks** — built on top of `testing.T`.

**Work in progress**

## Features

* Test **suites** with `BeforeAll`, `AfterEach`, etc.
* Custom test context `*T` with **plugin system** 🧩
* Extend and override test behavior
* Add or rename tests dynamically
* Run suite tests in parallel

## Quick Start

Define a suite and run it:

```go
type T struct { *testman.T } // define your own T

type MySuite struct{}

func (s MySuite) TestHello(t *T) {
	t.Log("Hello from test!")
}
```

```go
func Test(t *testing.T) {
	testman.Suite[MySuite, *T](t)
}
```

## Plugins

### Example: Plugin hooks

```go
type HelloLogger struct{ *testman.T }

func (h HelloLogger) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Hooks: plugin.Hooks{
			BeforeEach: func() {
				h.Log("👋 Hello from plugin!")
			},
		},
	}
}
```

Use it in your `T`:

```go
type T struct {
	*testman.T
	HelloLogger
}
```

And your suite:

```go
type MySuite struct{}

func (s MySuite) TestGreet(t *T) {
	t.Log("Running test")
}
```

### Example: Custom assertion helper

```go
type Assertions struct{ *testman.T }

func (a Assertions) RequireEqual(want, got any) {
	if want != got {
		a.Fatal("not equal")
	}
}
```

Usage:

```go
type T struct {
	*testman.T
	Assertions
}

type Suite struct{}

func (s Suite) TestCheck(t *T) {
	t.RequireEqual("hello", "hello")
}
```

### Example: Add a virtual test from plugin

```go
type AddExtraTest struct{}

func (AddExtraTest) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Plan: plugin.Plan{
			Add: func() []plugin.Test {
				return []plugin.Test{
					{
						Name: "extra test from plugin",
						Run: func(t plugin.T) {
							t.Log("👻 I was added dynamically!")
						},
					},
				}
			},
		},
	}
}
```

Use it in your `T` as usual.
