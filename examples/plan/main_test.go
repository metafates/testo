package main

import (
	"testing"

	"github.com/metafates/tego"
)

// Define your custom T and embed plugins
// along the base T provided by tego.
type T struct {
	*tego.T

	// These plugins modify tests in different ways.
	PluginWhichReversesTestOrder
	PluginWhichSkipsRandomTests
	PluginWhichOverridesLog

	Timer
}

// Entry point.
func Test(t *testing.T) {
	tego.RunSuite[Suite, *T](t)
}

type Suite struct{}

// Tests must have "Test" prefix.
func (Suite) Test1(t *T) {
	t.Parallel()

	t.Log("Hello from Test1")
}

func (Suite) Test2(t *T) {
	t.Parallel()

	t.Log("Hello from Test2")
}

// Cases funcs provide values for parametrized tests.
// They must have prefix "Cases".
func (Suite) CasesFruit() []string {
	return []string{"Apple", "Pear", "Grapefruit"}
}

func (Suite) CasesColor() []string {
	return []string{"Red", "Green", "Yellow"}
}

// Parametrized tests do not differ from other tests except
// that they accept second argument of type struct.
//
// Struct fields are matched against "CasesXXX" funcs.
// They are executed with all possible combinations of provided values.
func (Suite) Test3(t *T, params struct{ Fruit, Color string }) {
	t.Parallel()

	t.Logf("Hello from Test3. Fruit: %q, Color: %q", params.Fruit, params.Color)
}
