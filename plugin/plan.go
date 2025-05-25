package plugin

import (
	"testing"
)

type Plan struct {
	// Sort sorts the tests.
	//
	// It will not receive subtests as they can not be sorted
	// nor be known before running the parent tests.
	Sort func(a, b string) int

	// Add adds additional tests to be run.
	Add func() []Test
}

type T interface {
	testing.TB

	Run(name string, f func(t *testing.T)) bool
}

type Test struct {
	Name string
	Run  func(t T)
}

func mergePlans(plugins ...Plugin) Plan {
	return Plan{
		Sort: func(a, b string) int {
			// do not sort by default
			f := func(_, _ string) int { return 0 }

			for _, p := range plugins {
				if p.Plan.Sort != nil {
					f = p.Plan.Sort
				}
			}

			return f(a, b)
		},
		Add: func() []Test {
			var added []Test

			for _, p := range plugins {
				if p.Plan.Add != nil {
					added = append(added, p.Plan.Add()...)
				}
			}

			return added
		},
	}
}
