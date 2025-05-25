package main

import "testman"

// DOES NOT WORK YET, just an idea

type T struct{ *testman.T }

type Suite struct{}

func (Suite) CasesX() []int {
	return []int{1, 2, 3, 4}
}

func (Suite) CasesY() []string {
	return []string{"foo", "bar"}
}

func (Suite) TestFizz(t *T, args struct {
	X int
	Y string
},
) {
	t.Log(args.X, args.Y)
}
