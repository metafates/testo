package plugin

import (
	"github.com/metafates/testo/internal/pragma"
)

// Plan for running the tests.
type Plan struct {
	// Modify may filter or re-order planned tests in-place.
	// Nil values are ignored.
	Modify func(tests *[]PlannedTest)
}

type PlannedTest interface {
	pragma.DoNotImplement

	// Name of the test
	Name() string

	// TODO: other useful information about tests
}

func mergePlans(plugins ...Spec) Plan {
	return Plan{
		Modify: func(tests *[]PlannedTest) {
			for _, p := range plugins {
				// TODO: stop iterating if len(tests) == 0?

				if p.Plan.Modify != nil {
					p.Plan.Modify(tests)
				}
			}
		},
	}
}
