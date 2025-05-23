package testman

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type T interface {
	testing.TB

	Run(name string, f func(t T))
	Parallel()
}

type baseT struct {
	*testing.T
}

func (tt *baseT) Run(name string, f func(t T)) {
	tt.T.Run(name, func(t *testing.T) {
		wrapper := baseT{T: t}

		f(&wrapper)
	})
}

func Run[P, S any](t *testing.T, suite S, plugins ...Plugin) {
	tests := collectTests[S, P]()

	if len(tests) == 0 {
		return
	}

	_, err := composePlugins[P](plugins...)
	if err != nil {
		panic(err)
	}

	for _, handle := range tests {
		switch {
		case handle.Func != nil:
			wrapper := baseT{T: t}

			wrapper.Run(handle.Name, func(t T) { handle.Func(suite, t) })

		case handle.FuncP != nil:
			wrapperP := tp[P]{T: &baseT{T: t}}

			wrapperP.RunP(handle.Name, func(t TP[P]) { handle.FuncP(suite, t) })
		}
	}
}

func composePlugins[P any](plugins ...Plugin) (P, error) {
	pType := reflect.TypeOf((*P)(nil)).Elem()
	if pType.Kind() != reflect.Interface {
		var zero P

		return zero, fmt.Errorf("P (%s) must be an interface", pType)
	}

	wantMethods := make([]reflect.Method, 0, len(plugins))
	for i := range pType.NumMethod() {
		m := pType.Method(i)

		wantMethods = append(wantMethods, m)
	}

	haveMethods := make([]reflect.Value, 0, len(wantMethods))

	for _, m := range wantMethods {
		var found bool

		for _, p := range plugins {
			haveMethod := reflect.ValueOf(p).MethodByName(m.Name)

			if !haveMethod.IsZero() {
				haveMethods = append(haveMethods, haveMethod)
				found = true
				break
			}
		}

		if !found {
			var zero P

			return zero, fmt.Errorf("method not found: %s", m.Name)
		}
	}

	// TODO: compose these methods into value that would implement P
	panic("todo")
}

type testHandle[S, P any] struct {
	Name  string
	Func  func(S, T)
	FuncP func(S, TP[P])
}

func collectTests[S, P any]() []testHandle[S, P] {
	vt := reflect.TypeFor[S]()

	tests := make([]testHandle[S, P], 0, vt.NumMethod())

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		if !method.IsExported() {
			continue
		}

		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}

		switch f := method.Func.Interface().(type) {
		case func(S, T):
			tests = append(tests, testHandle[S, P]{
				Name: method.Name,
				Func: f,
			})

		case func(S, TP[P]):
			tests = append(tests, testHandle[S, P]{
				Name:  method.Name,
				FuncP: f,
			})
		}
	}

	return tests
}
