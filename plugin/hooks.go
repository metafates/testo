package plugin

import (
	"cmp"
	"slices"
)

type HookPriority int

const (
	// TryFirst indicates that this hook should be run as early as possible.
	TryFirst HookPriority = -1

	// TryLast indicates that this hook should be run as late as possible.
	TryLast HookPriority = 1
)

type Hook struct {
	// Priority defines execution order.
	// Lower values indicate that this hook should be run earlier than others and vice versa.
	// Zero value won't affect the order - it uses stable sort internally.
	//
	// See [TryFirst] and [TryLast] for predefined priority constants.
	Priority HookPriority

	// Func to be run for this hook.
	Func func()
}

// Run this hook.
// Calling this function with nil [Hook.Func] is safe and no op.
func (h Hook) Run() {
	if h.Func != nil {
		h.Func()
	}
}

// compare hooks based on their priority.
func (h Hook) compare(other Hook) int {
	return cmp.Compare(h.Priority, other.Priority)
}

type Hooks struct {
	// BeforeAll is called before all tests once.
	BeforeAll Hook

	// BeforeEach is called before each test and subtest.
	BeforeEach Hook

	// AfterEach is called after each test and subtest.
	AfterEach Hook

	// AfterAll is called after all tests once.
	AfterAll Hook
}

func mergeHooks(plugins ...Plugin) Hooks {
	beforeAll := make([]Hook, 0, len(plugins))
	beforeEach := make([]Hook, 0, len(plugins))
	afterEach := make([]Hook, 0, len(plugins))
	afterAll := make([]Hook, 0, len(plugins))

	for _, p := range plugins {
		if h := p.Hooks.BeforeAll; h.Func != nil {
			beforeAll = append(beforeAll, h)
		}

		if h := p.Hooks.BeforeEach; h.Func != nil {
			beforeEach = append(beforeEach, h)
		}

		if h := p.Hooks.AfterEach; h.Func != nil {
			afterEach = append(afterEach, h)
		}

		if h := p.Hooks.AfterAll; h.Func != nil {
			afterAll = append(afterAll, h)
		}
	}

	slices.SortStableFunc(beforeAll, Hook.compare)
	slices.SortStableFunc(beforeEach, Hook.compare)
	slices.SortStableFunc(afterEach, Hook.compare)
	slices.SortStableFunc(afterAll, Hook.compare)

	return Hooks{
		BeforeAll: Hook{
			Func: func() {
				for _, h := range beforeAll {
					h.Run()
				}
			},
		},
		BeforeEach: Hook{
			Func: func() {
				for _, h := range beforeEach {
					h.Run()
				}
			},
		},
		AfterEach: Hook{
			Func: func() {
				for _, h := range afterEach {
					h.Run()
				}
			},
		},
		AfterAll: Hook{
			Func: func() {
				for _, h := range afterAll {
					h.Run()
				}
			},
		},
	}
}
