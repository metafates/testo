package testman

import (
	"reflect"
	"strings"
	"testing"
)

type T struct {
	*testing.T
}

func New(t *testing.T) *T {
	return &T{T: t}
}

func (t *T) Run(name string, f func(t *T)) bool {
	return false
}

func Run[Suite any, T testing.TB](t T) {
	tests := collectSuiteTests[Suite, T]()

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	suite := *new(Suite)

	if i, ok := any(&suite).(beforeAller[T]); ok {
		i.BeforeAll(t)
	}

	for _, handle := range tests {
		suite := suite

		if i, ok := any(&suite).(beforeEacher[T]); ok {
			i.BeforeEach(t)
		}

		handle.Func(suite, t)

		if i, ok := any(&suite).(afterEacher[T]); ok {
			i.AfterEach(t)
		}
	}

	if i, ok := any(&suite).(afterAller[T]); ok {
		i.AfterAll(t)
	}
}

type testHandle[Suite any, T testing.TB] struct {
	Name string
	Func func(Suite, T)
}

func collectSuiteTests[Suite any, T testing.TB]() []testHandle[Suite, T] {
	vt := reflect.TypeFor[Suite]()

	tests := make([]testHandle[Suite, T], 0, vt.NumMethod())

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		if !method.IsExported() {
			continue
		}

		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}

		switch f := method.Func.Interface().(type) {
		case func(Suite, T):
			tests = append(tests, testHandle[Suite, T]{
				Name: method.Name,
				Func: f,
			})
		}
	}

	return tests
}
