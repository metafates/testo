package main

import (
	"testing"

	"github.com/metafates/tego"
)

type T = testing.T

// Also supported:
//
// type T = tego.T
// type T = struct{ *tego.T }

func Test(t *testing.T) {
	tego.RunSuite[Suite, *T](t)
}

type Suite struct{}

func (Suite) TestFoo(t *T) {
	t.Log("hi!")
	t.Fatal("oops")
}

func (Suite) CasesName() []string { return []string{"foo", "bar"} }

func (Suite) TestLen(t *T, params struct{ Name string }) {
	t.Log(params.Name)
}
