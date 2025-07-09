package allure

import (
	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

// Step is similar to [testo.Run], but if the step fails with fatal error,
// outer test execution will stop.
func Step[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) {
	t.Helper()

	var failure plugin.TestFailureKind

	fWrapper := func(t T) {
		t.Helper()

		defer func() {
			failure = testo.Inspect(t).FailureKind
		}()

		f(t)
	}

	if !testo.Run(t, name, fWrapper, options...) {
		// propagate fatal error
		if failure == plugin.TestFailureKindFatal {
			t.FailNow()
		}
	}
}

// Setup runs a [Step] marked as Setup in Allure report.
//
// You may want to use it in BeforeEach, BeforeAll hooks.
func Setup[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) {
	t.Helper()

	options = append(options, asSetup())

	Step(t, name, f, options...)
}

// TearDown runs a [Step] marked as TearDown in Allure report.
//
// You may want to use it in AfterEach, AfterAll hooks.
func TearDown[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) {
	t.Helper()

	options = append(options, asTearDown())

	Step(t, name, f, options...)
}
