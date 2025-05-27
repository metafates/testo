package foo

import (
	"time"

	"github.com/metafates/tego"
)

type Foo struct {
	*tego.T

	start time.Time
}

func (*Foo) New(t *tego.T) *Foo {
	return &Foo{T: t}
}

func (f *Foo) Measure() func() {
	f.Helper()

	f.start = time.Now()

	return func() {
		measure := time.Since(f.start)

		f.Log(f.Name() + " measured at " + measure.String())
	}
}

func (f *Foo) RequireTrue(value bool) {
	f.Helper()

	if !value {
		f.Fatalf("expected true, got false")
	}
}
