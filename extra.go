package testo

import "github.com/metafates/testo/plugin"

// ExtraInfo is extra information about T.
//
// Use [Inspect] to obtain it.
type ExtraInfo struct {
	// Plugins used by this T.
	Plugins []plugin.Plugin

	// Test holds information about current test.
	Test TestInfo

	// Panic is panic information.
	// It is nil if the test did not panic.
	Panic *PanicInfo
}

type PanicInfo struct {
	Value any
	Trace string
}

type TestInfo interface {
	isTestInfo()
}

type ParametrizedTestInfo struct {
	// BaseName of the test.
	BaseName string

	// Params passed for the current test case.
	Params map[string]any
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
