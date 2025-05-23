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
	tests := collectSuiteTests[Suite, T](t)

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	const (
		hookBeforeAll  = "BeforeAll"
		hookBeforeEach = "BeforeEach"

		hookAfterAll  = "AfterAll"
		hookAfterEach = "AfterEach"
	)

	tt := construct[T](&concreteT{T: t})

	var suite Suite

	callPluginHook(tt, hookBeforeAll)
	callSuiteHook(tt, &suite, hookBeforeAll)

	// if i, ok := any(&suite).(beforeAller[T]); ok {
	// i.BeforeAll(&tt)
	// }

	// so that AfterAll hooks will called after these tests even if they use Parallel().
	t.Run("X", func(t *testing.T) {
		for _, handle := range tests {
			suite := suite

			t.Run(handle.Name, func(t *testing.T) {
				tt := construct[T](&concreteT{T: t})

				callPluginHook(tt, hookBeforeEach)
				callSuiteHook(tt, &suite, hookBeforeEach)

				// if i, ok := any(&suite).(beforeEacher[T]); ok {
				// 	i.BeforeEach(&tt)
				// }

				handle.F(suite, &tt)

				callPluginHook(tt, hookAfterEach)
				callSuiteHook(tt, &suite, hookAfterEach)

				// if i, ok := any(&suite).(afterEacher[T]); ok {
				// 	i.AfterEach(&tt)
				// }
			})
		}
	})

	callPluginHook(tt, hookAfterAll)
	callSuiteHook(tt, &suite, hookAfterAll)

	// if i, ok := any(&suite).(afterAller[T]); ok {
	// i.AfterAll(&tt)
	// }
}

func Subtest[T constraint.T](t *T, name string, f func(t *T)) bool {
	// TODO: avoid dereferencing. With reflection?

	return (*t).Run(name, func(tt *testing.T) {
		t := construct[T](&concreteT{T: tt})

		f(&t)
	})
}

func callSuiteHook[T testing.TB](t T, suite any, name string) {
	sValue := reflect.ValueOf(suite)

	method := sValue.MethodByName(name)

	if method.IsValid() {
		f, ok := method.Interface().(func(*T))
		if !ok {
			t.Fatalf(
				"wrong signature for %[1]T.%[2]s, must be: func %[1]T.%[2]s(*%s)",
				suite, name, reflect.TypeFor[T](),
			)
		}

		f(&t)
	}
}

func callPluginHook[T testing.TB](t T, name string) {
	tValue := reflect.ValueOf(t)

	method := tValue.MethodByName(name)

	if method.IsValid() {
		f, ok := method.Interface().(func())
		if !ok {
			t.Fatalf(
				"wrong signature for %[1]T.%[2]s, must be: func %[1]T.%[2]s()",
				t, name,
			)
		}

		// f()
		_ = f
	}

	if tValue.Kind() != reflect.Struct {
		return
	}

	for i := range tValue.NumField() {
		field := tValue.Field(i)

		// TODO: make this recursive? (do we need this?)

		fieldMethod := field.MethodByName(name)

		if fieldMethod.IsValid() {
			f, ok := fieldMethod.Interface().(func())
			if !ok {
				t.Fatalf(
					"wrong signature for %[1]T.%[2]s, must be: func %[1]T.%[2]s()",
					t, name,
				)
			}

			f()
		}
	}
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

func collectSuiteTests[Suite any, T testing.TB](t *testing.T) []suiteTest[Suite, T] {
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

		default:
			t.Fatalf(
				"wrong signature for %[1]s.%[2]s, must be: func %[1]s.%[2]s(t *%s)",
				reflect.TypeFor[Suite](),
				method.Name,
				reflect.TypeFor[T](),
			)
		}
	}

	return tests
}
