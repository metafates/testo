package constraint

import "testing"

type T interface {
	testing.TB

	Run(name string, f func(t *testing.T)) bool
}
