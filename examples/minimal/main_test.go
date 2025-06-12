//go:build example

package main

import (
	"testing"

	"github.com/metafates/testo"
)

type T = testo.T

func Test(t *testing.T) {
	testo.RunSuite[Suite, *T](t)
}

type Suite struct{}

func (Suite) TestFoo(t *T) {
	if 2+2 != 4 {
		t.Fatal("2 + 2 must be equal 4")
	}
}
