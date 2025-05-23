package testman

import (
	"reflect"
	"strings"
	"testing"
)

type T struct {
	*testing.T
}

type customT[New testing.TB] interface {
	testing.TB

	New(*T) New
}

func (T) New(t *T) *T { return t }

func (t *T) Run(name string, f func(t *T)) bool {
	return false
}

func Run[Suite any, CT customT[CT]](t *testing.T) {
	tests := collectSuiteTests[Suite, CT]()

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	tt := (*new(CT)).New(&T{T: t})

	var suite Suite

	if i, ok := any(&suite).(beforeAller[CT]); ok {
		i.BeforeAll(&tt)
	}

	for _, handle := range tests {
		suite := suite

		t.Run(handle.Name, func(t *testing.T) {
			tt := tt.New(&T{T: t})

			if i, ok := any(&suite).(beforeEacher[CT]); ok {
				i.BeforeEach(&tt)
			}

			handle.F(suite, &tt)

			if i, ok := any(&suite).(afterEacher[CT]); ok {
				i.AfterEach(&tt)
			}
		})
	}

	if i, ok := any(&suite).(afterAller[CT]); ok {
		i.AfterAll(&tt)
	}
}

type suiteTest[Suite any, T testing.TB] struct {
	Name string
	F    func(Suite, *T)
}

func collectSuiteTests[Suite any, T testing.TB]() []suiteTest[Suite, T] {
	vt := reflect.TypeFor[Suite]()

	tests := make([]suiteTest[Suite, T], 0, vt.NumMethod())

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		if !method.IsExported() {
			continue
		}

		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}

		switch f := method.Func.Interface().(type) {
		case func(Suite, *T):
			tests = append(tests, suiteTest[Suite, T]{
				Name: method.Name,
				F:    f,
			})
		}
	}

	return tests
}
