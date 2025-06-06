package plugin

type MetaInfo struct {
	Plugins []Plugin
	Test    TestInfo
}

type TestInfo interface {
	isTestInfo()
}

type ParametrizedTestInfo struct {
	BaseName string
	Params   map[string]any
}

func (ParametrizedTestInfo) isTestInfo() {}

type RegularTestInfo struct{}

func (RegularTestInfo) isTestInfo() {}
