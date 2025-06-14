package plugin

// TInfo is extra information about T.
type TInfo struct {
	// Plugins used by this T.
	Plugins []Plugin

	// Test holds information about current test.
	Test TestInfo

	// Panic is panic information.
	// It is nil if the test did not panic.
	Panic *PanicInfo
}

// PanicInfo holds information for recovered panic.
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

// ParametrizedTestInfo is the information about parametrized test.
type ParametrizedTestInfo struct {
	// RawBaseName of the test.
	//
	// When defined like so:
	//
	//  func TestFoo(t T, params struct{ ... }) {}
	//
	// Testo will create a separate test for each case
	// and name it "TestFoo_case_N" where N is the case number.
	// Therefore, t.BaseName() would also equal to "TestFoo_case_N",
	// while this field would store "TestFoo".
	RawBaseName string

	// Params passed for the current test case.
	Params map[string]any
}

func (ParametrizedTestInfo) isTestInfo() {}

// RegularTestInfo is the information about regular (non-parametrized) test.
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
