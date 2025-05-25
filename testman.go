package testman

import (
	"reflect"
	"slices"
	"testing"

	"testman/internal/reflectutil"
	"testman/plugin"
)

// TODO: use real suite name
const wrapperTestName = "Suite"

// Suite will run the tests under the given suite.
//
// It also accepts options for the plugins which can be used to configure those plugins.
// See [plugin.Option].
func Suite[Suite any, T commonT](t *testing.T, options ...plugin.Option) {
	tt := construct[T](&concreteT{T: t}, nil, options...)
	plug := plugin.Merge(plugin.Collect(tt)...)
	tt.unwrap().overrides = plug.Overrides

	tests := collectSuiteTests[Suite, T](t)

	for _, a := range plug.Plan.Add() {
		tests = append(tests, suiteTest[Suite, T]{
			Name: a.Name,
			Run: func(_ Suite, t T) {
				a.Run(t)
			},
		})
	}

	for i := range tests {
		tests[i].Name = plug.Plan.Rename(tests[i].Name)
	}

	slices.SortFunc(tests, func(a, b suiteTest[Suite, T]) int {
		return plug.Plan.Sort(a.Name, b.Name)
	})

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	suiteHooks := collectSuiteHooks[Suite](tt)

	var suite Suite

	plug.Hooks.BeforeAll()
	defer plug.Hooks.AfterAll()

	suiteHooks.BeforeAll(suite, tt)
	defer suiteHooks.AfterAll(suite, tt)

	// so that AfterAll hooks will be called after these tests even if they use Parallel().
	t.Run(wrapperTestName, func(t *testing.T) {
		for _, test := range tests {
			var suiteClone Suite

			if s, ok := any(suite).(Cloner[Suite]); ok {
				suiteClone = s.Clone()
			} else {
				suiteClone = suite
			}

			t.Run(test.Name, func(t *testing.T) {
				subT := construct(&concreteT{T: t}, &tt)
				subPlug := plugin.Merge(plugin.Collect(subT)...)

				subT.unwrap().overrides = subPlug.Overrides

				subPlug.Hooks.BeforeEach()
				defer subPlug.Hooks.AfterEach()

				suiteHooks.BeforeEach(suiteClone, tt)
				defer suiteHooks.AfterEach(suiteClone, tt)

				test.Run(suite, subT)
			})
		}
	})
}

func Run[T commonT](t T, name string, f func(t T)) bool {
	plug := plugin.Merge(plugin.Collect(construct(t.unwrap(), &t))...)

	return t.Run(plug.Plan.Rename(name), func(tt *testing.T) {
		subT := construct(&concreteT{T: tt}, &t)

		plug := plugin.Merge(plugin.Collect(subT)...)

		plug.Hooks.BeforeEach()
		defer plug.Hooks.AfterEach()

		f(subT)
	})
}

func construct[V any](t *T, parent *V, options ...plugin.Option) V {
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
		options...,
	)

	return v
}

// initValue initializes value plugins (fields) with the given options
func initValue(t *T, value, parent reflect.Value, options ...plugin.Option) {
	const methodName = "New"

	var (
		constructor reflect.Value
		isPromoted  bool
	)

	if parent.IsValid() {
		constructor = parent.MethodByName(methodName)
		isPromoted = reflectutil.IsPromotedMethod(parent.Type(), methodName)
	} else {
		constructor = value.MethodByName(methodName)
		isPromoted = reflectutil.IsPromotedMethod(value.Type(), methodName)
	}

	if constructor.IsValid() && !isPromoted {
		mType := constructor.Type()

		// we can't assert an interface like .Interface().(func(*T, ...Option) G)
		// because we don't know anything about G here during compile type.

		isValidOut := mType.NumOut() == 1 && mType.Out(0) == value.Type()
		isValidIn := mType.NumIn() == 2 && mType.In(0) == reflect.TypeOf(t)

		if !isValidIn || !isValidOut {
			t.Fatalf(
				"wrong signature for %[1]s.New, must be: func (%[1]s) New(%T, %s...) %[1]s",
				value.Type().String(), t, reflect.TypeFor[plugin.Option](),
			)
		}

		res := constructor.CallSlice([]reflect.Value{
			reflect.ValueOf(t),
			reflect.ValueOf(options),
		})[0]

		value.Set(res)

		return
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
				initValue(t, field, parent.Field(i), options...)
			} else {
				initValue(t, field, reflect.ValueOf(nil), options...)
			}
		}
	}
}
