package main

import (
	"math/rand/v2"
	"reflect"

	"testman"

	"testman/plugin"
)

type RandomizeOrder struct{}

func (RandomizeOrder) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Plan: plugin.Plan{
			Sort: func(_, _ string) int { return rand.IntN(3) - 1 },
		},
	}
}

type Assertions struct{ *testman.T }

func (a Assertions) RequireEqual(want, got any) {
	a.Helper()

	if !reflect.DeepEqual(want, got) {
		a.Fatalf("not equal:\nwant: %v\ngot:  %v", want, got)
	}
}
