package tego

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/metafates/tego/plugin"
)

type CommonT interface {
	testing.TB

	Run(name string, f func(*testing.T)) bool

	unwrap() *T
}

func Inspect[T CommonT](t T) plugin.TestInfo {
	return t.unwrap().info
}

type (
	T struct {
		*testing.T

		parent    *T
		suiteName string
		plugin    plugin.Plugin
		panicInfo *PanicInfo

		// levelOptions stores option passes for
		// current level through [Run] or [RunSuite].
		levelOptions []plugin.Option

		info plugin.TestInfo
	}

	actualT = T
)

type PanicInfo struct {
	Value any
	Trace string
}

// Parallel signals that this test is to be run in parallel with (and only with)
// other parallel tests. When a test is run multiple times due to use of
// -test.count or -test.cpu, multiple instances of a single test never run in
// parallel with each other.
func (t *T) Parallel() {
	t.Helper()

	// This restricts the following pattern
	//
	// func TestFoo(t *testing.T) {
	//     t.Parallel() // level 1, this is ok
	//     t.Run("...", func(t *testing.T) { t.Parallel() }) // level 2, this is not supported
	// }
	//
	// the reason for that is that we won't be able to run AfterEach hook otherwise,
	// because test function will return control flow and continue running in a
	// separate goroutine later, thus triggering AfterEach to early.
	//
	// We can use t.Cleanup(AfterEach) to solve this, but if
	// AfterEach would call t.Run (which is common enough) the whole test will panic,
	// because running t.Run inside cleanup is not permitted (which makes sense, but unfortunate in our case).
	if t.level() == 2 {
		// TODO: add link to documentation or something so that user won't be left with questions.
		t.Log("running Parallel() at this level is not supported")

		return
	}

	t.plugin.Overrides.Parallel.Call(t.T.Parallel)()
}

// Chdir calls os.Chdir(dir) and uses Cleanup to restore the current
// working directory to its original value after the test. On Unix, it
// also sets PWD environment variable for the duration of the test.
//
// Because Chdir affects the whole process, it cannot be used
// in parallel tests or tests with parallel ancestors.
func (t *T) Chdir(dir string) {
	t.Helper()

	t.plugin.Overrides.Chdir.Call(t.T.Chdir)(dir)
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

// Context returns a context that is canceled just before
// Cleanup-registered functions are called.
//
// Cleanup functions can wait for any resources
// that shut down on Context.Done before the test or benchmark completes.
func (t *T) Context() context.Context {
	t.Helper()

	return t.plugin.Overrides.Context.Call(t.T.Context)()
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

// BaseName returns the base name for the current test.
// For example, given test "Test/Foo/Bar/MySubtest" it will return "MySubtest".
func (t *T) BaseName() string {
	segments := strings.Split(t.Name(), "/")

	if len(segments) == 0 {
		return ""
	}

	return segments[len(segments)-1]
}

// SuiteName returns current suite name.
func (t *T) SuiteName() string {
	if t.suiteName == "" {
		if t.parent != nil {
			return t.parent.SuiteName()
		}
	}

	return t.suiteName
}

// Panicked reports whether the function has panicked.
func (t *T) Panicked() bool {
	return t.panicInfo != nil
}

// PanicInfo returns information about panic occurred during this test.
func (t *T) PanicInfo() (PanicInfo, bool) {
	if t.panicInfo != nil {
		return *t.panicInfo, true
	}

	return PanicInfo{}, false
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

func (t *T) unwrap() *T {
	return t
}

// level indicates how deep this t is.
// That is, it shows the number of parents it has, zero if none.
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
		options = append(options, parent.levelOptions...)
		parent = parent.parent
	}

	return options
}
