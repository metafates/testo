package testo

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/metafates/testo/plugin"
)

// CommonT is the interface common for all [T] derivates.
type CommonT interface {
	testing.TB

	unwrap() *T
}

type (
	// T is a wrapper for [testing.T].
	// This is a core entity in testo.
	T struct {
		*testing.T

		parent *T
		plugin plugin.Spec

		// levelOptions stores option passes for
		// current level through [Run] or [RunSuite].
		levelOptions []plugin.Option

		suiteName string

		// info information required for [Inspect].
		info plugin.TInfo
	}

	actualT = T
)

// Inspect returns meta information about given t.
//
// Note that all plugins and suite tests share
// the same pointer to the underlying [T].
func Inspect[T CommonT](t T) plugin.TInfo {
	return t.unwrap().info
}

// Run runs f as a subtest of t called name. It runs f in a separate goroutine
// and blocks until f returns or calls t.Parallel to become a parallel test.
// Run reports whether f succeeded (or at least did not fail before calling t.Parallel).
//
// Run may be called simultaneously from multiple goroutines, but all such calls
// must return before the outer test function for t returns.
//
// Note: consider using [Run] function instead if you want to preserve type of T.
func (t *T) Run(name string, f func(t *T)) bool {
	return Run(t, name, f)
}

// SuiteName returns current suite name.
func (t *T) SuiteName() string {
	if t.suiteName != "" {
		return t.suiteName
	}

	parent := t.parent

	for parent != nil {
		if parent.suiteName != "" {
			return parent.suiteName
		}

		parent = parent.parent
	}

	return ""
}

// Parallel signals that this test is to be run in parallel with (and only with)
// other parallel tests. When a test is run multiple times due to use of
// -test.count or -test.cpu, multiple instances of a single test never run in
// parallel with each other.
//
// Note that running this method in the second level is not supported and treated as no op.
//
//	func TestFoo(t *T) {
//	   t.Parallel() // level 1, this is ok
//
//	   t.Run("...", func(t *T) {
//			// level 2, this is not supported
//			t.Parallel()
//
//	  		t.Run("...", func(t *T) {
//				// level 3, supported
//				t.Parallel()
//			})
//		})
//	}
func (t *T) Parallel() {
	t.Helper()

	// This restricts the pattern seen in the example above.
	//
	// The reason for that is that we won't be able to run AfterEach hook otherwise,
	// because test function will return control flow and continue running in a
	// separate goroutine later, thus triggering AfterEach too early.
	//
	// We could t.Cleanup(AfterEach) to solve this, but if
	// AfterEach would call t.Run (which is common enough) the whole test will panic,
	// because running t.Run inside cleanup is not permitted (which makes sense, but unfortunate in our case).
	if t.level() == 2 {
		// TODO: add link to documentation or something so that user won't be left with questions.
		t.Log("WARN: running Parallel() at this level is not supported and will be ignored")

		return
	}

	t.plugin.Overrides.Parallel.Call(t.T.Parallel)()
}

// Setenv calls os.Setenv(key, value) and uses Cleanup to
// restore the environment variable to its original value
// after the test.
//
// Because Setenv affects the whole process, it cannot be used
// in parallel tests or tests with parallel ancestors.
func (t *T) Setenv(key, value string) {
	t.Helper()

	t.plugin.Overrides.Setenv.Call(t.T.Setenv)(key, value)
}

// TempDir returns a temporary directory for the test to use.
// The directory is automatically removed when the test and
// all its subtests complete.
// Each subsequent call to t.TempDir returns a unique directory;
// if the directory creation fails, TempDir terminates the test by calling Fatal.
func (t *T) TempDir() string {
	t.Helper()

	return t.plugin.Overrides.TempDir.Call(t.T.TempDir)()
}

// Log formats its arguments using default formatting, analogous to Println,
// and records the text in the error log. For tests, the text will be printed only if
// the test fails or the -test.v flag is set. For benchmarks, the text is always
// printed to avoid having performance depend on the value of the -test.v flag.
func (t *T) Log(args ...any) {
	t.Helper()

	t.plugin.Overrides.Log.Call(t.T.Log)(args...)
}

// Logf formats its arguments according to the format, analogous to Printf, and
// records the text in the error log. A final newline is added if not provided. For
// tests, the text will be printed only if the test fails or the -test.v flag is
// set. For benchmarks, the text is always printed to avoid having performance
// depend on the value of the -test.v flag.
func (t *T) Logf(format string, args ...any) {
	t.Helper()

	t.plugin.Overrides.Logf.Call(t.T.Logf)(format, args...)
}

// Deadline reports the time at which the test binary will have
// exceeded the timeout specified by the -timeout flag.
//
// The ok result is false if the -timeout flag indicates “no timeout” (0).
func (t *T) Deadline() (time.Time, bool) {
	t.Helper()

	return t.plugin.Overrides.Deadline.Call(t.T.Deadline)()
}

// Errorf is equivalent to Logf followed by Fail.
func (t *T) Errorf(format string, args ...any) {
	t.Helper()

	t.plugin.Overrides.Errorf.Call(t.T.Errorf)(format, args...)
}

// Error is equivalent to Log followed by Fail.
func (t *T) Error(args ...any) {
	t.Helper()

	t.plugin.Overrides.Error.Call(t.T.Error)(args...)
}

// Skip is equivalent to Log followed by SkipNow.
func (t *T) Skip(args ...any) {
	t.Helper()

	t.plugin.Overrides.Skip.Call(t.T.Skip)(args...)
}

// SkipNow marks the test as having been skipped and stops its execution
// by calling [runtime.Goexit].
// If a test fails (see Error, Errorf, Fail) and is then skipped,
// it is still considered to have failed.
// Execution will continue at the next test or benchmark. See also FailNow.
// SkipNow must be called from the goroutine running the test, not from
// other goroutines created during the test. Calling SkipNow does not stop
// those other goroutines.
func (t *T) SkipNow() {
	t.Helper()

	t.plugin.Overrides.SkipNow.Call(t.T.SkipNow)()
}

// Skipf is equivalent to Logf followed by SkipNow.
func (t *T) Skipf(format string, args ...any) {
	t.Helper()

	t.plugin.Overrides.Skipf.Call(t.T.Skipf)(format, args...)
}

// Skipped reports whether the test was skipped.
func (t *T) Skipped() bool {
	t.Helper()

	return t.plugin.Overrides.Skipped.Call(t.T.Skipped)()
}

// Fail marks the function as having failed but continues execution.
func (t *T) Fail() {
	t.Helper()

	t.plugin.Overrides.Fail.Call(t.T.Fail)()
}

// FailNow marks the function as having failed and stops its execution
// by calling runtime.Goexit (which then runs all deferred calls in the
// current goroutine).
// Execution will continue at the next test or benchmark.
// FailNow must be called from the goroutine running the
// test or benchmark function, not from other goroutines
// created during the test. Calling FailNow does not stop
// those other goroutines.
func (t *T) FailNow() {
	t.Helper()

	t.plugin.Overrides.FailNow.Call(t.T.FailNow)()
}

// Failed reports whether the function has failed.
func (t *T) Failed() bool {
	t.Helper()

	return t.plugin.Overrides.Failed.Call(t.T.Failed)()
}

// Fatal is equivalent to Log followed by FailNow.
func (t *T) Fatal(args ...any) {
	t.Helper()

	t.plugin.Overrides.Fatal.Call(t.T.Fatal)(args...)
}

// Fatalf is equivalent to Logf followed by FailNow.
func (t *T) Fatalf(format string, args ...any) {
	t.Helper()

	t.plugin.Overrides.Fatalf.Call(t.T.Fatalf)(format, args...)
}

// Panicked reports whether the function has panicked.
func (t *T) Panicked() bool {
	return t.info.Panic != nil
}

// Name returns the name of the running (sub-) test or benchmark.
//
// The name will include the name of the test along with the names of
// any nested sub-tests. If two sibling sub-tests have the same name,
// Name will append a suffix to guarantee the returned name is unique.
func (t *T) Name() string {
	t.Helper()

	return t.plugin.Overrides.Name.Call(t.name)()
}

// name returns test name without [parallelWrapperTest] segment.
func (t *T) name() string {
	const sep = "/"

	name := t.T.Name()

	// segments in test name are always separate by forward slash /
	segments := strings.Split(name, sep)

	idx := slices.IndexFunc(segments, func(s string) bool {
		return s == parallelWrapperTest
	})
	if idx == -1 {
		return name
	}

	segments = slices.Delete(segments, idx, idx+1)

	name = strings.Join(segments, sep)

	return name
}

// unwrap the underlying T.
//
// It works since T's are embedded in user-defined structs.
func (t *T) unwrap() *T {
	return t
}

// level indicates how deep this t is.
// That is, it shows the number of parents it has and zero if none.
func (t *T) level() int {
	var level int

	parent := t.parent

	for parent != nil {
		level++
		parent = parent.parent
	}

	return level
}

func (t *T) options() []plugin.Option {
	options := t.levelOptions

	parent := t.parent

	for parent != nil {
		for _, o := range parent.levelOptions {
			if o.Propagate {
				options = append(options, o)
			}
		}

		parent = parent.parent
	}

	return options
}
