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

type (
	FuncLog  func(args ...any)
	FuncLogf func(format string, args ...any)
)

type Overrides struct {
	Name     func() string
	Parallel func()
	Chdir    func(dir string)
	Setenv   func(key string, value string)
	TempDir  func() string

	Log  Override[FuncLog]
	Logf Override[FuncLogf]

	Context  func() context.Context
	Deadline func() (deadline time.Time, ok bool)

	Errorf func(format string, args ...any)
	Error  func(args ...any)

	Skip    func(args ...any)
	SkipNow func()
	Skipf   func(format string, args ...any)
	Skipped func() bool

	Fail    func()
	FailNow func()
	Failed  func() bool

	Fatal  func(args ...any)
	Fatalf func(format string, args ...any)
}

// func (o *Overrides) Merge(other Overrides) {
// 	// TODO: combine multiple non nil like middlewares
//
// 	if other.Name != nil {
// 		o.Name = other.Name
// 	}
//
// 	if other.Parallel != nil {
// 		o.Parallel = other.Parallel
// 	}
//
// 	if other.Chdir != nil {
// 		o.Chdir = other.Chdir
// 	}
//
// 	if other.Setenv != nil {
// 		o.Setenv = other.Setenv
// 	}
//
// 	if other.TempDir != nil {
// 		o.TempDir = other.TempDir
// 	}
//
// 	if other.Log != nil {
// 		o.Log = other.Log
// 	}
//
// 	if other.Logf != nil {
// 		o.Logf = other.Logf
// 	}
//
// 	if other.Context != nil {
// 		o.Context = other.Context
// 	}
//
// 	if other.Deadline != nil {
// 		o.Deadline = other.Deadline
// 	}
//
// 	if other.Errorf != nil {
// 		o.Errorf = other.Errorf
// 	}
//
// 	if other.Error != nil {
// 		o.Error = other.Error
// 	}
//
// 	if other.Skip != nil {
// 		o.Skip = other.Skip
// 	}
//
// 	if other.SkipNow != nil {
// 		o.SkipNow = other.SkipNow
// 	}
//
// 	if other.Skipf != nil {
// 		o.Skipf = other.Skipf
// 	}
//
// 	if other.Skipped != nil {
// 		o.Skipped = other.Skipped
// 	}
//
// 	if other.Fail != nil {
// 		o.Fail = other.Fail
// 	}
//
// 	if other.FailNow != nil {
// 		o.FailNow = other.FailNow
// 	}
//
// 	if other.Failed != nil {
// 		o.Failed = other.Failed
// 	}
//
// 	if other.Fatal != nil {
// 		o.Fatal = other.Fatal
// 	}
//
// 	if other.Fatalf != nil {
// 		o.Fatalf = other.Fatalf
// 	}
// }

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
