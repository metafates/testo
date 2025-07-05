package allure

import (
	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

// Step is similar to [testo.Run], but will the step fails,
// the test execution will stop.
func Step[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) {
	if !testo.Run(t, name, f, options...) {
		t.FailNow()
	}
}

// Setup runs a subtest marked as Setup in Allure report.
//
// You may want to use in BeforeEach, BeforeAll hooks.
func Setup[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	options = append(options, asSetup())

	return testo.Run(t, name, f, options...)
}

// TearDown runs a subtest marked as TearDown in Allure report.
//
// You may want to use in AfterEach, AfterAll hooks.
func TearDown[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	options = append(options, asTearDown())

	return testo.Run(t, name, f, options...)
}
