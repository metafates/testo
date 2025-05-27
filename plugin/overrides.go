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

	FuncLog  func(args ...any)
	FuncLogf func(format string, args ...any)

	FuncContext  func() context.Context
	FuncDeadline func() (deadline time.Time, ok bool)

	FuncErrorf func(format string, args ...any)
	FuncError  func(args ...any)

	FuncSkip    func(args ...any)
	FuncSkipNow func()
	FuncSkipf   func(format string, args ...any)
	FuncSkipped func() bool

	FuncFail    func()
	FuncFailNow func()
	FuncFailed  func() bool

	FuncFatal  func(args ...any)
	FuncFatalf func(format string, args ...any)
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

func mergeOverrides(plugins ...Plugin) Overrides {
	return Overrides{
		Log: func(f FuncLog) FuncLog {
			for _, p := range plugins {
				if p.Overrides.Log != nil {
					f = p.Overrides.Log(f)
				}
			}

			return f
		},
		Logf: func(f FuncLogf) FuncLogf {
			for _, p := range plugins {
				o := p.Overrides.Logf

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Name: func(f FuncName) FuncName {
			for _, p := range plugins {
				o := p.Overrides.Name

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Parallel: func(f FuncParallel) FuncParallel {
			for _, p := range plugins {
				o := p.Overrides.Parallel

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Chdir: func(f FuncChdir) FuncChdir {
			for _, p := range plugins {
				o := p.Overrides.Chdir

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Setenv: func(f FuncSetenv) FuncSetenv {
			for _, p := range plugins {
				o := p.Overrides.Setenv

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		TempDir: func(f FuncTempDir) FuncTempDir {
			for _, p := range plugins {
				o := p.Overrides.TempDir

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Context: func(f FuncContext) FuncContext {
			for _, p := range plugins {
				o := p.Overrides.Context

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Deadline: func(f FuncDeadline) FuncDeadline {
			for _, p := range plugins {
				o := p.Overrides.Deadline

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Errorf: func(f FuncErrorf) FuncErrorf {
			for _, p := range plugins {
				o := p.Overrides.Errorf

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Error: func(f FuncError) FuncError {
			for _, p := range plugins {
				o := p.Overrides.Error

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Skip: func(f FuncSkip) FuncSkip {
			for _, p := range plugins {
				o := p.Overrides.Skip

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		SkipNow: func(f FuncSkipNow) FuncSkipNow {
			for _, p := range plugins {
				o := p.Overrides.SkipNow

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Skipf: func(f FuncSkipf) FuncSkipf {
			for _, p := range plugins {
				o := p.Overrides.Skipf

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Skipped: func(f FuncSkipped) FuncSkipped {
			for _, p := range plugins {
				o := p.Overrides.Skipped

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Fail: func(f FuncFail) FuncFail {
			for _, p := range plugins {
				o := p.Overrides.Fail

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		FailNow: func(f FuncFailNow) FuncFailNow {
			for _, p := range plugins {
				o := p.Overrides.FailNow

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Failed: func(f FuncFailed) FuncFailed {
			for _, p := range plugins {
				o := p.Overrides.Failed

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Fatal: func(f FuncFatal) FuncFatal {
			for _, p := range plugins {
				o := p.Overrides.Fatal

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
		Fatalf: func(f FuncFatalf) FuncFatalf {
			for _, p := range plugins {
				o := p.Overrides.Fatalf

				if o != nil {
					f = o(f)
				}
			}

			return f
		},
	}
}
