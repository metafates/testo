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
		t *testing.T

		overrides plugin.Overrides
	}

	concreteT = T
)

func (*T) New(t *T) *T { return t }

func (t *T) Name() string {
	if t.overrides.Name != nil {
		return t.overrides.Name()
	}

	return t.name()
}

func (t *T) name() string {
	name := t.t.Name()

	idx := strings.Index(name, wrapperTestName)

	if idx >= 0 {
		// +1 for slash
		return name[idx+len(wrapperTestName)+1:]
	}

	return name
}

func (t *T) Parallel() {
	if t.overrides.Parallel != nil {
		t.overrides.Parallel()
		return
	}

	t.t.Parallel()
}

func (t *T) Chdir(dir string) {
	if t.overrides.Chdir != nil {
		t.overrides.Chdir(dir)
		return
	}

	t.t.Chdir(dir)
}

func (t *T) Setenv(key string, value string) {
	if t.overrides.Setenv != nil {
		t.overrides.Setenv(key, value)
		return
	}

	t.t.Setenv(key, value)
}

func (t *T) TempDir() string {
	if t.overrides.TempDir != nil {
		return t.overrides.TempDir()
	}

	return t.t.TempDir()
}

func (t *T) Log(args ...any) {
	t.overrides.Log(t.t.Log)(args...)
}

func (t *T) Logf(format string, args ...any) {
	if t.overrides.Logf == nil {
		t.t.Logf(format, args...)
		return
	}

	t.overrides.Logf(t.t.Logf)(format, args...)
}

func (t *T) Context() context.Context {
	if t.overrides.Context != nil {
		return t.overrides.Context()
	}

	return t.t.Context()
}

func (t *T) Deadline() (deadline time.Time, ok bool) {
	if t.overrides.Deadline != nil {
		return t.overrides.Deadline()
	}

	return t.t.Deadline()
}

func (t *T) Errorf(format string, args ...any) {
	if t.overrides.Errorf != nil {
		t.overrides.Errorf(format, args...)
		return
	}

	t.t.Errorf(format, args...)
}

func (t *T) Error(args ...any) {
	if t.overrides.Error != nil {
		t.overrides.Error(args...)
		return
	}

	t.t.Error(args...)
}

func (t *T) Skip(args ...any) {
	if t.overrides.Skip != nil {
		t.overrides.Skip(args...)
		return
	}

	t.t.Skip(args...)
}

func (t *T) SkipNow() {
	if t.overrides.SkipNow != nil {
		t.overrides.SkipNow()
		return
	}

	t.t.SkipNow()
}

func (t *T) Skipf(format string, args ...any) {
	if t.overrides.Skipf != nil {
		t.overrides.Skipf(format, args...)
		return
	}

	t.t.Skipf(format, args...)
}

func (t *T) Skipped() bool {
	if t.overrides.Skipped != nil {
		return t.overrides.Skipped()
	}

	return t.t.Skipped()
}

func (t *T) Fail() {
	if t.overrides.Fail != nil {
		t.overrides.Fail()
		return
	}

	t.t.Fail()
}

func (t *T) FailNow() {
	if t.overrides.FailNow != nil {
		t.overrides.FailNow()
		return
	}

	t.t.FailNow()
}

func (t *T) Failed() bool {
	if t.overrides.Failed != nil {
		return t.overrides.Failed()
	}

	return t.t.Failed()
}

func (t *T) Fatal(args ...any) {
	if t.overrides.Fatal != nil {
		t.overrides.Fatal(args...)
		return
	}

	t.t.Fatal(args...)
}

func (t *T) Fatalf(format string, args ...any) {
	if t.overrides.Fatalf != nil {
		t.overrides.Fatalf(format, args...)
		return
	}

	t.t.Fatalf(format, args...)
}

func (t *T) BaseName() string {
	segments := strings.Split(t.Name(), "/")

	if len(segments) == 0 {
		return ""
	}

	return segments[len(segments)-1]
}

func (t *T) Cleanup(f func()) {
	t.t.Cleanup(f)
}

func (t *T) Helper() {
	t.t.Helper()
}

func (t *T) Run(name string, f func(t *testing.T)) bool {
	return t.t.Run(name, f)
}

func (t *T) T() *testing.T {
	return t.t
}

func (t *T) TT() *T {
	return t
}
