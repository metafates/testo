package testman

import "testing"

type commonT interface {
	runner
	fataller

	TT() *T
}

type runner interface {
	Run(string, func(*testing.T)) bool
}

type fataller interface {
	Fatalf(string, ...any)
}
