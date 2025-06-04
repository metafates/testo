package plugin

import (
	"slices"
	"testing"
)

// Plan for running the tests.
type Plan struct {
	Modify func(tests []PlannedTest) []PlannedTest
}

type T interface {
	testing.TB

	Run(name string, f func(t *testing.T)) bool
}

type PlannedTest interface {
	GetName() string

	// TODO: other useful information about tests
}

func mergePlans(plugins ...PluginSpec) Plan {
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
