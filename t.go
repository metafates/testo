package testman

import (
	"context"
	"strings"
	"testing"
	"time"

	"testman/plugin"
)

type (
	T struct {
		*testing.T

		overrides plugin.Overrides
	}

	concreteT = T
)

func (*T) New(t *T) *T { return t }

func (t *T) Name() string {
	t.Helper()

	return t.overrides.Name.Call(t.name)()
}

func (t *T) name() string {
	name := t.T.Name()

	idx := strings.Index(name, wrapperTestName)

	if idx >= 0 {
		// +1 for slash
		return name[idx+len(wrapperTestName)+1:]
	}

	return name
}

func (t *T) Parallel() {
	t.Helper()

	t.overrides.Parallel.Call(t.T.Parallel)()
}

func (t *T) Chdir(dir string) {
	t.Helper()

	t.overrides.Chdir.Call(t.T.Chdir)(dir)
}

func (t *T) Setenv(key, value string) {
	t.Helper()

	t.overrides.Setenv.Call(t.T.Setenv)(key, value)
}

func (t *T) TempDir() string {
	t.Helper()

	return t.overrides.TempDir.Call(t.T.TempDir)()
}

func (t *T) Log(args ...any) {
	t.Helper()

	t.overrides.Log.Call(t.T.Log)(args...)
}

func (t *T) Logf(format string, args ...any) {
	t.Helper()

	t.overrides.Logf.Call(t.T.Logf)(format, args...)
}

func (t *T) Context() context.Context {
	t.Helper()

	return t.overrides.Context.Call(t.T.Context)()
}

func (t *T) Deadline() (deadline time.Time, ok bool) {
	t.Helper()

	return t.overrides.Deadline.Call(t.T.Deadline)()
}

func (t *T) Errorf(format string, args ...any) {
	t.Helper()

	t.overrides.Errorf.Call(t.T.Errorf)(format, args...)
}

func (t *T) Error(args ...any) {
	t.Helper()

	t.overrides.Error.Call(t.T.Error)(args...)
}

func (t *T) Skip(args ...any) {
	t.Helper()

	t.overrides.Skip.Call(t.T.Skip)(args...)
}

func (t *T) SkipNow() {
	t.Helper()

	t.overrides.SkipNow.Call(t.T.SkipNow)()
}

func (t *T) Skipf(format string, args ...any) {
	t.Helper()

	t.overrides.Skipf.Call(t.T.Skipf)(format, args...)
}

func (t *T) Skipped() bool {
	t.Helper()

	return t.overrides.Skipped.Call(t.T.Skipped)()
}

func (t *T) Fail() {
	t.Helper()

	t.overrides.Fail.Call(t.T.Fail)()
}

func (t *T) FailNow() {
	t.Helper()

	t.overrides.FailNow.Call(t.T.FailNow)()
}

func (t *T) Failed() bool {
	t.Helper()

	return t.overrides.Failed.Call(t.T.Failed)()
}

func (t *T) Fatal(args ...any) {
	t.Helper()

	t.overrides.Fatal.Call(t.T.Fatal)(args...)
}

func (t *T) Fatalf(format string, args ...any) {
	t.Helper()

	t.overrides.Fatalf.Call(t.T.Fatalf)(format, args...)
}

func (t *T) BaseName() string {
	segments := strings.Split(t.Name(), "/")

	if len(segments) == 0 {
		return ""
	}

	return segments[len(segments)-1]
}

func (t *T) unwrap() *T {
	return t
}
