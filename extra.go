package testo

import "github.com/metafates/testo/plugin"

type ExtraInfo struct {
	parent    func() ExtraInfo
	suiteName string

	Plugins []plugin.Plugin
	Test    TestInfo
	Panic   *PanicInfo
}

func (e ExtraInfo) SuiteName() string {
	if e.suiteName != "" {
		return e.suiteName
	}

	if e.parent == nil {
		return ""
	}

	return e.parent().SuiteName()
}

func (e ExtraInfo) Parent() (ExtraInfo, bool) {
	if e.parent == nil {
		return ExtraInfo{}, false
	}

	return e.parent(), true
}

type PanicInfo struct {
	Value any
	Trace string
}

type TestInfo interface {
	isTestInfo()
}

type ParametrizedTestInfo struct {
	BaseName string
	Params   map[string]any
}

func (ParametrizedTestInfo) isTestInfo() {}

type RegularTestInfo struct {
	// RawBaseName is the raw "unformatted" base name of this test.
	//
	// For example:
	//
	//   Run(t, "my test name!?", func(...) { ... })
	//
	// t.BaseName() would equal to my_test_name,
	// while this field would equal to "my test name!?".
	//
	// This only applies to subtests, for regular test
	// functions BaseName() and this field are equal.
	RawBaseName string
}

func (RegularTestInfo) isTestInfo() {}
