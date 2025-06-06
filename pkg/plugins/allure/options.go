package allure

import (
	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

type option func(*Allure)

func WithCategories(categories ...Category) plugin.Option {
	return plugin.Option{
		Value: option(func(a *Allure) {
			a.categories = categories
		}),
	}
}

func WithOutputPath(path string) plugin.Option {
	return plugin.Option{
		Value: option(func(a *Allure) {
			a.outputPath = path
		}),
	}
}

func withTitle(title string) plugin.Option {
	return plugin.Option{
		Value: option(func(a *Allure) {
			a.titleOverwrite = title
		}),
	}
}

func asSetup() plugin.Option {
	return plugin.Option{
		Value: option(func(a *Allure) {
			a.stage = stageSetup
		}),
	}
}

func asTearDown() plugin.Option {
	return plugin.Option{
		Value: option(func(a *Allure) {
			a.stage = stageTearDown
		}),
	}
}

// TearDown runs a subtest which will be marked as Setup in Allure report.
func Setup[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	options = append(options, asSetup(), withTitle(name))

	return testo.Run(t, name, f, options...)
}

// TearDown runs a subtest which will be marked as TearDown in Allure report.
func TearDown[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	options = append(options, asTearDown(), withTitle(name))

	return testo.Run(t, name, f, options...)
}
