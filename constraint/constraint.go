package constraint

import "testing"

type T interface {
	testing.TB

	Run(name string, f func(*testing.T)) bool
}
