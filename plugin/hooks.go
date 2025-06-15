package plugin

import (
	"cmp"
	"slices"
)

// HookPriority is the [Hook] priority.
// It defines when a hook should be invoked when other hooks are available.
type HookPriority int

const (
	// TryFirst indicates that this hook should be run as early as possible.
	TryFirst HookPriority = -1

	// TryLast indicates that this hook should be run as late as possible.
	TryLast HookPriority = 1
)

// Hook is the plugin hook.
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

// Hooks defines all hooks a plugin can define.
type Hooks struct {
	// BeforeAll is called before all tests once.
	BeforeAll Hook

	// BeforeEach is called before each test.
	BeforeEach Hook

	// BeforeEachSub is called before each subtest.
	BeforeEachSub Hook

	// AfterEachSub is called after each subtest.
	AfterEachSub Hook

	// AfterEach is called after each test.
	AfterEach Hook

	// AfterAll is called after all tests once.
	AfterAll Hook
}

func mergeHooks(plugins ...Spec) Hooks {
	beforeAll := make([]Hook, 0, len(plugins))
	beforeEach := make([]Hook, 0, len(plugins))
	beforeEachSub := make([]Hook, 0, len(plugins))
	afterEachSub := make([]Hook, 0, len(plugins))
	afterEach := make([]Hook, 0, len(plugins))
	afterAll := make([]Hook, 0, len(plugins))

	for _, p := range plugins {
		if h := p.Hooks.BeforeAll; h.Func != nil {
			beforeAll = append(beforeAll, h)
		}

		if h := p.Hooks.BeforeEach; h.Func != nil {
			beforeEach = append(beforeEach, h)
		}

		if h := p.Hooks.BeforeEachSub; h.Func != nil {
			beforeEachSub = append(beforeEachSub, h)
		}

		if h := p.Hooks.AfterEachSub; h.Func != nil {
			afterEachSub = append(afterEachSub, h)
		}

		if h := p.Hooks.AfterEach; h.Func != nil {
			afterEach = append(afterEach, h)
		}

		if h := p.Hooks.AfterAll; h.Func != nil {
			afterAll = append(afterAll, h)
		}
	}

	run := func(hooks []Hook) func() {
		slices.SortStableFunc(hooks, Hook.compare)

		return func() {
			for _, h := range hooks {
				h.Run()
			}
		}
	}

	return Hooks{
		BeforeAll:     Hook{Func: run(beforeAll)},
		BeforeEach:    Hook{Func: run(beforeEach)},
		BeforeEachSub: Hook{Func: run(beforeEachSub)},
		AfterEachSub:  Hook{Func: run(afterEachSub)},
		AfterEach:     Hook{Func: run(afterEach)},
		AfterAll:      Hook{Func: run(afterAll)},
	}
}
