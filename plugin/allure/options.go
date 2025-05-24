package allure

type Option func(*Allure)

func WithOutputPath(path string) Option {
	return func(a *Allure) {
		a.outputPath = path
	}
}
