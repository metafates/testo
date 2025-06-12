//go:build example

package main

import (
	"reflect"

	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

type RandomizeOrder struct{}

func (RandomizeOrder) Plugin() plugin.Spec {
	return plugin.Spec{
		Plan: plugin.Plan{},
	}
}

type Assertions struct{ *testo.T }

func (a Assertions) RequireEqual(want, got any) {
	a.Helper()

	if !reflect.DeepEqual(want, got) {
		a.Fatalf("not equal:\nwant: %v\ngot:  %v", want, got)
	}
}
