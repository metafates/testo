package testman

import (
	"reflect"
	"strings"
	"testing"

	"testman/internal/constraint"
)

type T struct {
	*testing.T
}

type customT[New testing.TB] interface {
	testing.TB

	New(*T) New
}

func (*T) New(t *T) *T { return t }

func Run[Suite any, CT testing.TB](t *testing.T) {
	tests := collectSuiteTests[Suite, CT]()

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	// tt := (*new(CT)).New(&T{T: t})

	tt := construct[CT](&T{T: t})

	var suite Suite

	if i, ok := any(&suite).(beforeAller[CT]); ok {
		i.BeforeAll(&tt)
	}

	for _, handle := range tests {
		suite := suite

		t.Run(handle.Name, func(t *testing.T) {
			tt := construct[CT](&T{T: t})

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

func Subtest[CT constraint.T](t *CT, name string, f func(t *CT)) bool {
	return (*t).Run(name, func(tt *testing.T) {
		t := construct[CT](&T{T: tt})

		f(&t)
	})
}

func construct[V any](t *T) V {
	return rConstruct(t, reflect.ValueOf(new(V))).Interface().(V)
}

func rConstruct(t *T, value reflect.Value) reflect.Value {
	methodNew := value.MethodByName("New")

	if methodNew.IsValid() {
		mType := methodNew.Type()

		isValidIn := mType.NumIn() == 1 && mType.In(0) == reflect.TypeOf(t)
		isValidOut := mType.NumOut() == 1 && mType.Out(0) == value.Type()

		if isValidIn && isValidOut {
			return methodNew.Call([]reflect.Value{reflect.ValueOf(t)})[0]
		}
	}

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return value
	}

	for i := range value.NumField() {
		field := value.Field(i)

		if field.CanSet() {
			field.Set(rConstruct(t, field))
		}
	}

	return value
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
