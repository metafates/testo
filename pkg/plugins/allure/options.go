package allure

import (
	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

type Option func(*Allure)

func WithOutputPath(path string) plugin.Option {
	return plugin.Option{
		Value: Option(func(a *Allure) {
			a.outputPath = path
		}),
	}
}

func asSetup() plugin.Option {
	return plugin.Option{
		Value: Option(func(a *Allure) {
			a.stage = stageSetup
		}),
	}
}

func asTearDown() plugin.Option {
	return plugin.Option{
		Value: Option(func(a *Allure) {
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
	options = append(options, asSetup())

	return testo.Run(t, name, f, options...)
}

// TearDown runs a subtest which will be marked as TearDown in Allure report.
func TearDown[T testo.CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	options = append(options, asTearDown())

	return testo.Run(t, name, f, options...)
}
