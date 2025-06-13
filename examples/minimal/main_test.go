//go:build example

package main

import (
	"testing"

	"github.com/metafates/testo"
)

type T = *testo.T

func Test(t *testing.T) {
	testo.RunSuite[Suite, T](t)
}

type Suite struct{}

func (Suite) BeforeEach(t T) {
	t.Logf("Starting: %s", t.Name())
}

func (Suite) AfterEach(t T) {
	t.Logf("Finished: %s", t.Name())
}

func (Suite) TestAdd(t T) {
	t.Log(t.Name())
	if 2+2 != 4 {
		t.Fatal("2 + 2 must be equal 4")
	}
}

func (Suite) TestAddButParametrized(t T, params struct{ A, B int }) {
	if params.A+params.B != params.B+params.A {
		t.Errorf("%[1]d + %[2]d != %[2]d + %[1]d", params.A, params.B)
	}
}

func (Suite) CasesA() []int {
	return []int{1, 2, 3, 4, 5}
}

func (Suite) CasesB() []int {
	return []int{11, 1000, 13}
}
