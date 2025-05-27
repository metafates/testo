package tego

import (
	"testing"
)

// Cloner can clone itself.
type Cloner[Self any] interface {
	// Clone returns a new instance cloned from the caller.
	Clone() Self
}

type commonT interface {
	testing.TB

	runner

	unwrap() *T
}

type runner interface {
	Run(name string, f func(*testing.T)) bool
}

type fataller interface {
	Fatalf(format string, args ...any)
}
