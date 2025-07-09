package plugin

import (
	"time"
)

// TODO: check which functions we really need to allow for override. E.g. do we need [Overrides.Parallel]?
// TODO: support overriding t methods from the future versions (e.g. Context and Chdir).

// Overrides defines all builtin methods of T a plugin can override.
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
	Setenv   Override[FuncSetenv]
	TempDir  Override[FuncTempDir]
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
	// FuncName describes [testing.T.Name] method.
	FuncName func() string

	// FuncParallel describes [testing.T.Parallel] method.
	FuncParallel func()

	// FuncSetenv describes [testing.T.Setenv] method.
	FuncSetenv func(key, value string)

	// FuncTempDir describes [testing.T.TempDir] method.
	FuncTempDir func() string

	// FuncLog describes [testing.T.Log] method.
	FuncLog func(args ...any)

	// FuncLogf describes [testing.T.Logf] method.
	FuncLogf func(format string, args ...any)

	// FuncDeadline describes [testing.T.Deadline] method.
	FuncDeadline func() (deadline time.Time, ok bool)

	// FuncErrorf describes [testing.T.Errorf] method.
	FuncErrorf func(format string, args ...any)

	// FuncError describes [testing.T.Error] method.
	FuncError func(args ...any)

	// FuncSkip describes [testing.T.Skip] method.
	FuncSkip func(args ...any)

	// FuncSkipNow describes [testing.T.SkipNow] method.
	FuncSkipNow func()

	// FuncSkipf describes [testing.T.Skipf] method.
	FuncSkipf func(format string, args ...any)

	// FuncSkipped describes [testing.T.Skipped] method.
	FuncSkipped func() bool

	// FuncFail describes [testing.T.Fail] method.
	FuncFail func()

	// FuncFailNow describes [testing.T.FailNow] method.
	FuncFailNow func()

	// FuncFailed describes [testing.T.Failed] method.
	FuncFailed func() bool

	// FuncFatal describes [testing.T.Fatal] method.
	FuncFatal func(args ...any)

	// FuncFatalf describes [testing.T.Fatalf] method.
	FuncFatalf func(format string, args ...any)
)

// Override for the function.
//
// Nil value is valid and represents absence of override.
type Override[F any] func(f F) F

// Call returns an overridden f.
// If override is nil f is returned as is.
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

func mergeOverride[F any](
	plugins []Spec,
	getter func(Overrides) Override[F],
) func(F) F {
	return func(f F) F {
		for _, p := range plugins {
			if o := getter(p.Overrides); o != nil {
				f = o(f)
			}
		}

		return f
	}
}
