package plugin

import (
	"context"
	"time"
)

// TODO: check which functions we really need to allow for override. E.g. do we need [Overrides.Parallel]?

type Overrides struct {
	// Log overrides [testing.T.Log] function.
	//
	// Note that overriding this function won't override internal log used for Logf
	// or other functions such as Error, Skip, Fatal, etc.
	Log Override[FuncLog]

	// Logf overrides [testing.T.Logf] function.
	//
	// Note that overriding this function won't override internal log used for Logf
	// or other functions such as Errorf, Skipf, Fatalf, etc.
	Logf Override[FuncLogf]

	Name     Override[FuncName]
	Parallel Override[FuncParallel]
	Chdir    Override[FuncChdir]
	Setenv   Override[FuncSetenv]
	TempDir  Override[FuncTempDir]
	Context  Override[FuncContext]
	Deadline Override[FuncDeadline]
	Errorf   Override[FuncErrorf]
	Error    Override[FuncError]
	Skip     Override[FuncSkip]
	SkipNow  Override[FuncSkipNow]
	Skipf    Override[FuncSkipf]
	Skipped  Override[FuncSkipped]
	Fail     Override[FuncFail]
	FailNow  Override[FuncFailNow]
	Failed   Override[FuncFailed]
	Fatal    Override[FuncFatal]
	Fatalf   Override[FuncFatalf]
}

type (
	FuncName     func() string
	FuncParallel func()
	FuncChdir    func(dir string)
	FuncSetenv   func(key, value string)
	FuncTempDir  func() string
	FuncLog      func(args ...any)
	FuncLogf     func(format string, args ...any)
	FuncContext  func() context.Context
	FuncDeadline func() (deadline time.Time, ok bool)
	FuncErrorf   func(format string, args ...any)
	FuncError    func(args ...any)
	FuncSkip     func(args ...any)
	FuncSkipNow  func()
	FuncSkipf    func(format string, args ...any)
	FuncSkipped  func() bool
	FuncFail     func()
	FuncFailNow  func()
	FuncFailed   func() bool
	FuncFatal    func(args ...any)
	FuncFatalf   func(format string, args ...any)
)

// Override for the function.
//
// Nil value is valid and represents absence of override.
type Override[F any] func(f F) F

func (o Override[F]) Call(f F) F {
	if o == nil {
		return f
	}

	return o(f)
}

//nolint:funlen // splitting this into subfunctons would make it worse
func mergeOverrides(plugins ...Spec) Overrides {
	return Overrides{
		Log: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncLog] {
				return o.Log
			},
		),
		Logf: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncLogf] {
				return o.Logf
			},
		),
		Name: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncName] {
				return o.Name
			},
		),
		Parallel: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncParallel] {
				return o.Parallel
			},
		),
		Chdir: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncChdir] {
				return o.Chdir
			},
		),
		Setenv: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncSetenv] {
				return o.Setenv
			},
		),
		TempDir: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncTempDir] {
				return o.TempDir
			},
		),
		Context: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncContext] {
				return o.Context
			},
		),
		Deadline: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncDeadline] {
				return o.Deadline
			},
		),
		Errorf: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncErrorf] {
				return o.Errorf
			},
		),
		Error: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncError] {
				return o.Error
			},
		),
		Skip: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncSkip] {
				return o.Skip
			},
		),
		SkipNow: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncSkipNow] {
				return o.SkipNow
			},
		),
		Skipf: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncSkipf] {
				return o.Skipf
			},
		),
		Skipped: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncSkipped] {
				return o.Skipped
			},
		),
		Fail: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncFail] {
				return o.Fail
			},
		),
		FailNow: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncFailNow] {
				return o.FailNow
			},
		),
		Failed: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncFailed] {
				return o.Failed
			},
		),
		Fatal: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncFatal] {
				return o.Fatal
			},
		),
		Fatalf: mergeOverride(
			plugins,
			func(o Overrides) Override[FuncFatalf] {
				return o.Fatalf
			},
		),
	}
}

func mergeOverride[Fn any](
	plugins []Spec,
	getter func(Overrides) Override[Fn],
) func(Fn) Fn {
	return func(f Fn) Fn {
		for _, p := range plugins {
			if o := getter(p.Overrides); o != nil {
				f = o(f)
			}
		}

		return f
	}
}
