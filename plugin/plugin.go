package plugin

import (
	"context"
	"reflect"
	"time"

	"testman/internal/reflectutil"
)

type Plugin struct {
	Hooks     Hooks
	Overrides Overrides
}

type Hooks struct {
	BeforeAll  func()
	BeforeEach func()
	AfterEach  func()
	AfterAll   func()
}

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
					if p.Overrides.Logf != nil {
						f = p.Overrides.Logf(f)
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
	rootPlugin := scanPlugin(v, visited)

	rv := reflectutil.Elem(reflect.ValueOf(v))

	plugins := []Plugin{rootPlugin}

	if rv.Kind() != reflect.Struct {
		return []Plugin{rootPlugin}
	}

	for i := range rv.NumField() {
		field := rv.Field(i)

		if field.IsValid() && rv.Type().Field(i).IsExported() {
			plugins = append(plugins, collectPlugins(field.Interface(), visited)...)
		}
	}

	return plugins
}

func scanPlugin(v any, visited map[uintptr]struct{}) Plugin {
	var p Plugin

	if v, ok := v.(interface{ Hooks() Hooks }); ok {
		ptr := reflect.ValueOf(v.Hooks).Pointer()

		if _, ok := visited[ptr]; !ok {
			p.Hooks = v.Hooks()
		}

		visited[ptr] = struct{}{}
	}

	if v, ok := v.(interface{ Overrides() Overrides }); ok {
		ptr := reflect.ValueOf(v.Overrides).Pointer()

		if _, ok := visited[ptr]; !ok {
			p.Overrides = v.Overrides()
		}

		visited[ptr] = struct{}{}
	}

	return p
}
