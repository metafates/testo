package main

import (
	"reflect"

	"github.com/metafates/tego"
	"github.com/metafates/tego/plugin"
)

type RandomizeOrder struct{}

func (RandomizeOrder) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Plan: plugin.Plan{},
	}
}

type Assertions struct{ *tego.T }

func (a Assertions) RequireEqual(want, got any) {
	a.Helper()

	if !reflect.DeepEqual(want, got) {
		a.Fatalf("not equal:\nwant: %v\ngot:  %v", want, got)
	}
}
