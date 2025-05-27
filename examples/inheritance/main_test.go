package main

import (
	"testing"

	"github.com/metafates/tego"
)

// This is the base T which contains common plugins and maybe add it's own methods
type BaseT struct {
	*tego.T

	RandomizeOrder
}

// This method will be available for [T] because it embeds [BaseT]
func (t BaseT) Hi() {
	t.Helper()
	t.Log("Hi")
}

type T struct {
	// that's how we can "inherit" [BaseT].
	BaseT

	// even though [BaseT] already includes [tego.T]
	// it is required to embed it here once again,
	// otherwise we will get "ambiguous selector" compile error (go type system quirks, don't bother).
	*tego.T

	// Here we include another plugin on top of ones inherited from the [BaseT].
	Assertions
}

func Test(t *testing.T) {
	tego.RunSuite[Suite, *T](t)
}

type Suite struct{}

func (Suite) TestFoo(t *T) {
	// this comes from [BaseT]
	t.Hi()

	t.Log("Hi from TestFoo!")
}

func (Suite) TestBar(t *T) {
	t.Log("Hi from TestBar!")

	// this comes from [Assertions] plugin
	t.RequireEqual("hello", "oops")
}

func (Suite) TestFizz(t *T) {
	t.Log("Hi from TestFizz!")
}
