package main

import (
	"testing"

	"github.com/metafates/testo"
)

type T struct{ *testo.T }

func Test(t *testing.T) {
	testo.RunSuite[*Suite, *T](t)
}

type Suite struct{}

func (Suite) AfterAll(t *T) {
	t.Log(">>>>>> suite tear down", t.Name())
}

func (Suite) AfterEach(t *T) {
	t.Log(">>>>>> single test tear down", t.Name())
}

func (Suite) CasesX() []int {
	return []int{1, 2, 3, 4}
}

func (Suite) CasesY() []string {
	return []string{"foo", "bar"}
}

func (Suite) CasesZ() []bool { return []bool{true, false} }

func (Suite) TestFizz(t *T, params struct {
	X int
	Y string
	Z bool
},
) {
	t.Parallel()

	// filter out invalid combinations
	if params.Z && params.X%2 == 0 {
		t.Skip()
	}

	t.Log(params.X, params.Y, params.Z)
}

func (Suite) TestBuzz(t *T, params struct{ X int }) {
	t.Log("hi!", params.X)
}

// You can also use structs

type AddCase struct{ A, B, Want int }

func (Suite) CasesAdd() []AddCase {
	return []AddCase{
		{A: 2, B: 3, Want: 5},
		{A: 0, B: -3, Want: -3},
	}
}

func (Suite) TestAdditions(t *T, params struct{ Add AddCase }) {
	add := params.Add

	if add.A+add.B != add.Want {
		t.Error("%d + %d != %d", add.A, add.B, add.Want)
	}
}
