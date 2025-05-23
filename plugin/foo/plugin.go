package foo

import (
	"fmt"
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
	f.start = time.Now()

	return func() {
		measure := time.Since(f.start)

		fmt.Println(f.Name() + " measured at " + measure.String())
	}
}

func (f *Foo) RequireTrue(value bool) {
	f.Helper()

	if !value {
		f.Fatalf("expected true, got false")
	}
}
