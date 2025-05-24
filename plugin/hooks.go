package plugin

type Hooks struct {
	// BeforeAll is called before all tests once.
	BeforeAll func()

	// BeforeEach is called before each test and subtest.
	BeforeEach func()

	// AfterEach is called after each test and subtest.
	AfterEach func()

	// AfterAll is called after all tests once.
	AfterAll func()
}

func mergeHooks(plugins ...Plugin) Hooks {
	return Hooks{
		BeforeAll: func() {
			for _, p := range plugins {
				if p.Hooks.BeforeAll != nil {
					p.Hooks.BeforeAll()
				}
			}
		},
		BeforeEach: func() {
			for _, p := range plugins {
				if p.Hooks.BeforeEach != nil {
					p.Hooks.BeforeEach()
				}
			}
		},
		AfterEach: func() {
			for _, p := range plugins {
				if p.Hooks.AfterEach != nil {
					p.Hooks.AfterEach()
				}
			}
		},
		AfterAll: func() {
			for _, p := range plugins {
				if p.Hooks.AfterAll != nil {
					p.Hooks.AfterAll()
				}
			}
		},
	}
}
