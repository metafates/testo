package tego

import (
	"testing"
)

type commonT interface {
	testing.TB

	Run(name string, f func(*testing.T)) bool

	unwrap() *T
}
