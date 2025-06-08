package plugin

import (
	"slices"
)

// Plan for running the tests.
type Plan struct {
	Modify func(tests []PlannedTest) []PlannedTest
}

type PlannedTest interface {
	GetName() string

	// TODO: other useful information about tests
}

func mergePlans(plugins ...Spec) Plan {
	return Plan{
		Modify: func(tests []PlannedTest) []PlannedTest {
			for _, p := range plugins {
				// TODO: stop iterating if len(tests) == 0?

				if p.Plan.Modify != nil {
					tests = p.Plan.Modify(slices.Clone(tests))
				}
			}

			return tests
		},
	}
}
