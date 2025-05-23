package foo

import (
	"time"

	"testman"
)

type Foo struct {
	*testman.T

	start time.Time
}

func (Foo) New(t *testman.T) Foo {
	return Foo{T: t}
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
