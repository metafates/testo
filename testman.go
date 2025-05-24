package testman

import (
	"reflect"
	"strings"
	"testing"

	"testman/internal/constraint"
	"testman/plugin"
)

const wrapperTestName = "!"

const (
	hookBeforeAll  = "BeforeAll"
	hookBeforeEach = "BeforeEach"
	hookAfterAll   = "AfterAll"
	hookAfterEach  = "AfterEach"
)

func Suite[Suite any, T constraint.CommonT](t *testing.T) {
	tests := collectSuiteTests[Suite, T](t)

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	tt := construct[T](&concreteT{t: t}, nil)

	plug := plugin.Merge(plugin.Collect(tt)...)

	var suite Suite

	plug.Hooks.BeforeAll()
	callSuiteHook(tt, &suite, hookBeforeAll)

	// so that AfterAll hooks will be called after these tests even if they use Parallel().
	t.Run(wrapperTestName, func(t *testing.T) {
		for _, handle := range tests {
			suite := suite

			t.Run(handle.Name, func(t *testing.T) {
				subT := construct(&concreteT{t: t}, &tt)
				subPlug := plugin.Merge(plugin.Collect(tt)...)

				subPlug.Hooks.BeforeEach()
				callSuiteHook(subT, &suite, hookBeforeEach)

				defer subPlug.Hooks.AfterEach()
				defer callSuiteHook(subT, &suite, hookAfterEach)

				handle.F(suite, subT)
			})
		}
	})

	plug.Hooks.AfterAll()
	callSuiteHook(tt, &suite, hookAfterAll)
}

func Run[T constraint.CommonT](t T, name string, f func(t T)) bool {
	// TODO: avoid dereferencing. With reflection?

	return t.Run(name, func(tt *testing.T) {
		subT := construct(&concreteT{t: tt}, &t)

		callPluginHook(subT, hookBeforeEach)
		defer callPluginHook(subT, hookAfterEach)

		f(subT)
	})
}

func callSuiteHook[T constraint.Fataller](t T, suite any, name string) {
	sValue := elem(reflect.ValueOf(suite))

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

		isValidIn := mType.NumIn() == 1 && mType.In(0) == reflect.TypeOf(t)
		isValidOut := mType.NumOut() == 1 && mType.Out(0) == value.Type()

		if isValidIn && isValidOut {
			res := methodNew.Call([]reflect.Value{reflect.ValueOf(t)})[0]

			value.Set(res)

			return
		}
	}

	value = elem(value)
	parent = elem(parent)

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

func elem(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return v
}

type suiteTest[Suite any, T any] struct {
