package main

import (
	"testing"

	"testman"
)

// DOES NOT WORK YET, just an idea

type T struct{ *testman.T }

func Test(t *testing.T) {
	testman.Suite[Suite, *T](t)
}

type Suite struct{}

func (Suite) CasesX() []int {
	return []int{1, 2, 3, 4}
}

func (Suite) CasesY() []string {
	return []string{"foo", "bar"}
}

func (Suite) CasesZ() []bool { return []bool{true, false} }

func (Suite) TestFizz(t *T, args struct {
	X int
	Y string
	Z bool
},
) {
	// filter out invalid combinations
	if args.Z && args.X%2 == 0 {
		t.Skip()
	}

	t.Log(args.X, args.Y, args.Z)
}

func (Suite) TestBuzz(t *T) {
	t.Log("hi!")
}
