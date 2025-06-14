package allure

import (
	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

type option func(*Allure)

// WithCategories adds [custom categories] to the report.
//
// [custom categories]: https://allurereport.org/docs/categories/#custom-categories
func WithCategories(categories ...Category) plugin.Option {
	return plugin.Option{
		Value: option(func(a *Allure) {
			a.categories = append(a.categories, categories...)
		}),
	}
}

// WithOutputDir sets output directory for test results.
//
// By default, it is "current working directory/allure-results".
func WithOutputDir(dir string) plugin.Option {
	return plugin.Option{
		Value: option(func(a *Allure) {
			a.outputDir = dir
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

// Setup runs a subtest marked as Setup in Allure report.
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
func TearDown[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	options = append(options, asTearDown())

	return testo.Run(t, name, f, options...)
}
