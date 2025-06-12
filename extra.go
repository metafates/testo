package testo

import "github.com/metafates/testo/plugin"

// Inspect returns meta information about given t.
//
// Note that all plugins and suite tests share
// the same pointer to the underlying [T].
func Inspect[T CommonT](t T) ExtraInfo {
	return t.unwrap().extra
}

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
	// Value returned by recover().
	Value any

	// Trace is a stack trace for this panic.
	Trace string
}

// TestInfo is a enum which is
// either [ParametrizedTestInfo] or [RegularTestInfo].
type TestInfo interface {
	isTestInfo()
}

type ParametrizedTestInfo struct {
	// BaseName of the test.
	//
	// When defined like so:
	//
	//  func TestFoo(t T, params struct{ ... }) {}
	//
	// Testo will create a separate test for each case
	// and name it "TestFoo_case_N" where N is the case number.
	// Therefore, t.Name() would also equal to "TestFoo_case_N",
	// while this field would store "TestFoo".
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
	// while this field would equal to "my test name!?" (the same as passed).
	//
	// This only applies to subtests, for regular test
	// functions BaseName() and this field are equal.
	RawBaseName string

	// IsSubtest specifies whether this test is a (possibly nested) subtest.
	IsSubtest bool
}

func (RegularTestInfo) isTestInfo() {}
