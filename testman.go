package testman

import (
	"reflect"
	"strings"
	"testing"

	"testman/internal/reflectutil"
	"testman/plugin"
)

const wrapperTestName = "!"

const (
	hookBeforeAll  = "BeforeAll"
	hookBeforeEach = "BeforeEach"
	hookAfterAll   = "AfterAll"
	hookAfterEach  = "AfterEach"
)

func Suite[Suite any, T commonT](t *testing.T) {
	tests := collectSuiteTests[Suite, T](t)

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	tt := construct[T](&concreteT{T: t}, nil)
	plug := plugin.Merge(plugin.Collect(tt)...)

	var suite Suite

	plug.Hooks.BeforeAll()
	callSuiteHook(tt, &suite, hookBeforeAll)

	// so that AfterAll hooks will be called after these tests even if they use Parallel().
	t.Run(wrapperTestName, func(t *testing.T) {
		for _, handle := range tests {
			var suiteClone Suite

			if s, ok := any(suite).(Cloner[Suite]); ok {
				suiteClone = s.Clone()
			} else {
				suiteClone = suite
			}

			t.Run(handle.Name, func(t *testing.T) {
				subT := construct(&concreteT{T: t}, &tt)
				subPlug := plugin.Merge(plugin.Collect(subT)...)

				subT.unwrap().overrides = subPlug.Overrides

				subPlug.Hooks.BeforeEach()
				defer subPlug.Hooks.AfterEach()

				callSuiteHook(subT, &suiteClone, hookBeforeEach)
				defer callSuiteHook(subT, &suiteClone, hookAfterEach)

				handle.F(suite, subT)
			})
		}
	})

	plug.Hooks.AfterAll()
	callSuiteHook(tt, &suite, hookAfterAll)
}

func Run[T commonT](t T, name string, f func(t T)) bool {
	return t.Run(name, func(tt *testing.T) {
		subT := construct(&concreteT{T: tt}, &t)

		plug := plugin.Merge(plugin.Collect(subT)...)

		plug.Hooks.BeforeEach()
		defer plug.Hooks.AfterEach()

		f(subT)
	})
}

func callSuiteHook[T fataller](t T, suite any, name string) {
	sValue := reflectutil.Elem(reflect.ValueOf(suite))

	method := sValue.MethodByName(name)

	if method.IsValid() {
		f, ok := method.Interface().(func(T))
		if !ok {
			t.Fatalf(
				"wrong signature for %[1]T.%[2]s, must be: func %[1]T.%[2]s(*%s)",
				suite, name, reflect.TypeFor[T](),
			)
		}

		f(t)
	}
}

func construct[V any](t *T, parent *V) V {
	value := reflect.ValueOf(*new(V))

	if value.Kind() == reflect.Pointer && value.IsNil() {
		value = reflect.New(value.Type().Elem())
	}

	parentValue := reflect.ValueOf(parent)
	if parent != nil {
		parentValue = reflect.ValueOf(*parent)
	}

	v := value.Interface().(V)

	initValue(
		t,
		reflect.ValueOf(&v),
		parentValue,
	)

	return v
}

func initValue(t *T, value, parent reflect.Value) {
	var methodNew reflect.Value

	if parent.IsValid() {
		methodNew = parent.MethodByName("New")
	} else {
		methodNew = value.MethodByName("New")
	}

	if methodNew.IsValid() {
		mType := methodNew.Type()

		// we can't assert an interface like .Interface().(func(*T) G)
		// because we don't know anything about G here during compile type.

		isValidIn := mType.NumIn() == 1 && mType.In(0) == reflect.TypeOf(t)
		isValidOut := mType.NumOut() == 1 && mType.Out(0) == value.Type()

		if isValidIn && isValidOut {
			res := methodNew.Call([]reflect.Value{reflect.ValueOf(t)})[0]

			value.Set(res)

			return
		}
	}

	value = reflectutil.Elem(value)
	parent = reflectutil.Elem(parent)

	if value.Kind() != reflect.Struct {
		return
	}

	for i := range value.NumField() {
		field := value.Field(i)

		if field.CanSet() {
			if parent.IsValid() {
				initValue(t, field, parent.Field(i))
			} else {
				initValue(t, field, reflect.ValueOf(nil))
			}
		}
	}
}

type suiteTest[Suite any, T any] struct {
	Name string
	F    func(Suite, T)
}

func collectSuiteTests[Suite any, T fataller](t *testing.T) []suiteTest[Suite, T] {
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
		case func(Suite, T):
			tests = append(tests, suiteTest[Suite, T]{
				Name: method.Name,
				F:    f,
			})

		default:
			t.Fatalf(
				"wrong signature for %[1]s.%[2]s, must be: func %[1]s.%[2]s(t %s)",
				reflect.TypeFor[Suite](),
				method.Name,
				reflect.TypeFor[T](),
			)
		}
	}

	return tests
}
