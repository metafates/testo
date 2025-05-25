package plugin

import (
	"testing"
)

type Plan struct {
	// Sort tests.
	//
	// It will not receive subtests as they can not be sorted
	// nor be known before running the parent tests.
	Sort func(a, b string) int

	// Add additional tests to be run.
	Add func() []Test

	// Rename each test including subtests.
	Rename func(name string) string
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

			// TODO: merge sorting functions somehow.
			//
			// Also, it would be better to iterate backwards and break after first non-nil
			// but for the sake of simplicity let's leave it as is for now.
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
		Rename: func(name string) string {
			for _, p := range plugins {
				if p.Plan.Rename != nil {
					name = p.Plan.Rename(name)
				}
			}

			return name
		},
	}
}
