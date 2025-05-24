package plugin

import (
	"context"
	"reflect"
	"time"

	"testman/internal/reflectutil"
)

type Option struct {
	Value any
}

type Plugin struct {
	Hooks     Hooks
	Overrides Overrides
}

type Hooks struct {
	// BeforeAll is called before all tests once.
	BeforeAll func()

	// BeforeEach is called before each test and subtest.
	BeforeEach func()

	// AfterEach is called after each test and subtest.
	AfterEach func()

	// AfterAll is called after all tests once.
	AfterAll func()
}

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

type Overrides struct {
	Name     Override[FuncName]
	Parallel Override[FuncParallel]
	Chdir    Override[FuncChdir]
	Setenv   Override[FuncSetenv]
	TempDir  Override[FuncTempDir]

	Log  Override[FuncLog]
	Logf Override[FuncLogf]

	Context  Override[FuncContext]
	Deadline Override[FuncDeadline]

	Errorf Override[FuncErrorf]
	Error  Override[FuncError]

	Skip    Override[FuncSkip]
	SkipNow Override[FuncSkipNow]
	Skipf   Override[FuncSkipf]
	Skipped Override[FuncSkipped]

	Fail    Override[FuncFail]
	FailNow Override[FuncFailNow]
	Failed  Override[FuncFailed]

	Fatal  Override[FuncFatal]
	Fatalf Override[FuncFatalf]
}

func Merge(plugins ...Plugin) Plugin {
	return Plugin{
		Hooks: Hooks{
			BeforeAll: func() {
				for _, p := range plugins {
					if p.Hooks.BeforeAll != nil {
						p.Hooks.BeforeAll()
					}
				}
			},
			BeforeEach: func() {
				for _, p := range plugins {
					if p.Hooks.BeforeEach != nil {
						p.Hooks.BeforeEach()
					}
				}
			},
			AfterEach: func() {
				for _, p := range plugins {
					if p.Hooks.AfterEach != nil {
						p.Hooks.AfterEach()
					}
				}
			},
			AfterAll: func() {
				for _, p := range plugins {
					if p.Hooks.AfterAll != nil {
						p.Hooks.AfterAll()
					}
				}
			},
		},
		Overrides: Overrides{
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
		},
	}
}

func Collect(v any) []Plugin {
	return collectPlugins(v, make(map[uintptr]struct{}))
}

func collectPlugins(v any, visited map[uintptr]struct{}) []Plugin {
	var plugins []Plugin

	rootPlugin, ok := scanPlugin(v, visited)
	if ok {
		plugins = append(plugins, rootPlugin)
	}

	rv := reflectutil.Elem(reflect.ValueOf(v))

	if rv.Kind() == reflect.Struct {
		for i := range rv.NumField() {
			field := rv.Field(i)

			if field.IsValid() && rv.Type().Field(i).IsExported() {
				plugins = append(plugins, collectPlugins(field.Interface(), visited)...)
			}
		}
	}

	return plugins
}

func scanPlugin(v any, visited map[uintptr]struct{}) (Plugin, bool) {
	// We could (and will) access the same method twice if it was promoted
	// from an embed type, which will result calling it twice later.
	//
	// Therefore we maintain visited pointers set so that we won't collect the same method twice.

	if v, ok := v.(interface{ Plugin() Plugin }); ok {
		ptr := reflect.ValueOf(v.Plugin).Pointer()

		if _, ok := visited[ptr]; !ok {
			visited[ptr] = struct{}{}

			return v.Plugin(), true
		}
	}

	return Plugin{}, false
}
