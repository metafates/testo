package allure

import "github.com/metafates/tego/plugin"

type Option func(*Allure)

func WithOutputPath(path string) plugin.Option {
	return plugin.Option{
		Value: Option(func(a *Allure) {
			a.outputPath = path
		}),
	}
}

func AsSetup() plugin.Option {
	return plugin.Option{
		Value: Option(func(a *Allure) {
			a.stage = stageSetup
		}),
	}
}

func AsTearDown() plugin.Option {
	return plugin.Option{
		Value: Option(func(a *Allure) {
			a.stage = stageTearDown
		}),
	}
}
