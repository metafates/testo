package allure

import (
	"github.com/metafates/testo/plugin"
)

type option func(*Allure)

// WithCategories adds [custom categories] to the report.
// This option should be passed to the top-level [testo.RunSuite] call.
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
