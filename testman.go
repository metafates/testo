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

type concreteT = T

type customT[New testing.TB] interface {
	testing.TB

	New(*T) New
}

func (*T) New(t *T) *T { return t }

func Run[Suite any, T testing.TB](t *testing.T) {
	tests := collectSuiteTests[Suite, T]()

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	// tt := (*new(CT)).New(&T{T: t})

	tt := construct[T](&concreteT{T: t})

	var suite Suite

	if i, ok := any(&suite).(beforeAller[T]); ok {
		i.BeforeAll(&tt)
	}

	for _, handle := range tests {
		suite := suite

		t.Run(handle.Name, func(t *testing.T) {
			tt := construct[T](&concreteT{T: t})

			if i, ok := any(&suite).(beforeEacher[T]); ok {
				i.BeforeEach(&tt)
			}

			handle.F(suite, &tt)

			if i, ok := any(&suite).(afterEacher[T]); ok {
				i.AfterEach(&tt)
			}
		})
	}

	if i, ok := any(&suite).(afterAller[T]); ok {
		i.AfterAll(&tt)
	}
}

func Subtest[T constraint.T](t *T, name string, f func(t *T)) bool {
	return (*t).Run(name, func(tt *testing.T) {
		t := construct[T](&concreteT{T: tt})

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
