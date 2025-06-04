package main

import (
	"testing"

	"github.com/metafates/testo"
)

type T = testo.T

// Also supported:
//
// type T = struct{ *testo.T }

func Test(t *testing.T) {
	testo.RunSuite[Suite, *T](t)
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
