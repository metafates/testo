package constraint

import "testing"

type CommonT interface {
	Runner
	Fataller
}

type Runner interface {
	Run(string, func(*testing.T)) bool
}

type Fataller interface {
	Fatalf(string, ...any)
}
