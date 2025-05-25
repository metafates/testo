package testman

import (
	"reflect"
	"slices"
	"testing"

	"testman/internal/reflectutil"
	"testman/internal/stack"
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

	tests := applyPlan(plug.Plan, collectSuiteTests[Suite, T](t))

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

	// wrap all tests so that AfterAll hooks will
	// be called after these tests even if they use Parallel().
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
	value := reflectutil.Filled[V]()

	inits := stack.New[func()]()

	initValue(
		t,
		reflect.ValueOf(&value),
		reflect.ValueOf(parent),
		&inits,
		options...,
	)

	for {
		init, ok := inits.Pop()
		if !ok {
			break
		}

		init()
	}

	return value
}

func initValue(
	t *T,
	value, parent reflect.Value,
	inits *stack.Stack[func()],
	options ...plugin.Option,
) {
	// init T's
	if value.Type() == reflect.TypeOf(t) {
		value.Set(reflect.ValueOf(t))
		return
	}

	const methodName = "Init"

	if value.CanAddr() {
		value = value.Addr()
	}

	if parent.CanAddr() {
		parent = parent.Addr()
	}

	initFunc := value.MethodByName(methodName)
	isPromoted := reflectutil.IsPromotedMethod(value.Type(), methodName)

	if initFunc.IsValid() && !isPromoted {
		method := initFunc.Type()

		isValidOut := method.NumOut() == 0
		isValidIn := method.NumIn() == 2 && method.In(0) == parent.Type()

		if !isValidIn || !isValidOut {
			t.Fatalf(
				"wrong signature for %[1]T.Init, must be: func (%[1]T) Init(%[1]T, ...%s)",
				value.Interface(), reflect.TypeFor[plugin.Option](),
			)
		}

		parent := parent

		inits.Push(func() {
			initFunc.CallSlice([]reflect.Value{
				parent,
				reflect.ValueOf(options),
			})
		})
	}

	value = reflectutil.Elem(value)
	parent = reflectutil.Elem(parent)

	if value.Kind() != reflect.Struct {
		return
	}

	for i := range value.NumField() {
		field := value.Field(i)

		if !field.CanSet() {
			continue
		}

		if parent.IsValid() {
			initValue(t, field, parent.Field(i), inits, options...)
		} else {
			initValue(t, field, reflect.New(field.Type()), inits, options...)
		}
	}
}

func applyPlan[Suite any, T commonT](
	plan plugin.Plan,
	tests []suiteTest[Suite, T],
) []suiteTest[Suite, T] {
	for _, a := range plan.Add() {
		tests = append(tests, suiteTest[Suite, T]{
			Name: a.Name,
			Run: func(_ Suite, t T) {
				a.Run(t)
			},
		})
	}

	for i := range tests {
		tests[i].Name = plan.Rename(tests[i].Name)
	}

	slices.SortFunc(tests, func(a, b suiteTest[Suite, T]) int {
		return plan.Sort(a.Name, b.Name)
	})

	return tests
}
