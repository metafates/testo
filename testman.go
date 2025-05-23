package testman

import (
	"reflect"
	"strings"
	"testing"

	"testman/internal/constraint"
)

const wrapperTestName = "!"

type T struct {
	*testing.T
}

type concreteT = T

func (*T) New(t *T) *T { return t }

func (t *T) Name() string {
	name := t.T.Name()

	idx := strings.Index(name, wrapperTestName)

	if idx >= 0 {
		return name[idx+2:]
	}

	return name
}

func (t *T) BaseName() string {
	segments := strings.Split(t.Name(), "/")

	if len(segments) == 0 {
		return ""
	}

	return segments[len(segments)-1]
}

const (
	hookBeforeAll  = "BeforeAll"
	hookBeforeEach = "BeforeEach"

	hookAfterAll  = "AfterAll"
	hookAfterEach = "AfterEach"
)

func Suite[Suite any, T testing.TB](t *testing.T) {
	tests := collectSuiteTests[Suite, T](t)

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	tt := construct[T](&concreteT{T: t}, nil)

	var suite Suite

	callPluginHook(tt, hookBeforeAll)
	callSuiteHook(tt, &suite, hookBeforeAll)

	// so that AfterAll hooks will called after these tests even if they use Parallel().
	t.Run(wrapperTestName, func(t *testing.T) {
		for _, handle := range tests {
			suite := suite

			t.Run(handle.Name, func(t *testing.T) {
				subT := construct(&concreteT{T: t}, &tt)

				callPluginHook(subT, hookBeforeEach)
				callSuiteHook(subT, &suite, hookBeforeEach)

				defer callPluginHook(subT, hookAfterEach)
				defer callSuiteHook(subT, &suite, hookAfterEach)

				handle.F(suite, subT)
			})
		}
	})

	callPluginHook(tt, hookAfterAll)
	callSuiteHook(tt, &suite, hookAfterAll)
}

func Run[T constraint.T](t T, name string, f func(t T)) bool {
	// TODO: avoid dereferencing. With reflection?

	return t.Run(name, func(tt *testing.T) {
		subT := construct(&concreteT{T: tt}, &t)

		callPluginHook(subT, hookBeforeEach)
		defer callPluginHook(subT, hookAfterEach)

		f(subT)
	})
}

func callSuiteHook[T testing.TB](t T, suite any, name string) {
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

func callPluginHook[T testing.TB](t T, name string) {
	tValue := elem(reflect.ValueOf(t))

	if tValue.Kind() != reflect.Struct {
		return
	}

	for i := range tValue.NumField() {
		field := tValue.Field(i)

		// TODO: make this recursive? (do we need this?)

		method := field.MethodByName(name)

		if method.IsValid() {
			f, ok := method.Interface().(func())
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

type suiteTest[Suite any, T testing.TB] struct {
	Name string
	F    func(Suite, T)
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

// func inherit(parent, child reflect.Value) {
// 	mInherit := child.MethodByName("Inherit")
//
// 	if mInherit.IsValid() {
// 		mType := mInherit.Type()
//
// 		isValidIn := mType.NumIn() == 1 && mType.In(0) == parent.Type()
// 		isValidOut := mType.NumOut() == 0
//
// 		if isValidIn && isValidOut {
// 			mInherit.Call([]reflect.Value{parent})
// 			return
// 		}
// 	}
// }
