package plugin

import "testing"

type T interface {
	testing.TB

	Run(name string, f func(t *testing.T)) bool
}

type Test struct {
	Name string
	Run  func(t T)
}
