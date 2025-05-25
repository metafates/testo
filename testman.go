package testman

import (
	"fmt"
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

	tt.unwrap().plugin = plugin.Merge(plugin.Collect(tt)...)

	tests := applyPlan(tt.unwrap().plugin.Plan, collectSuiteTests[Suite, T](t))

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	suiteHooks := collectSuiteHooks[Suite](tt)

	var suite Suite

	tt.unwrap().plugin.Hooks.BeforeAll()
	defer tt.unwrap().plugin.Hooks.AfterAll()

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
				subT.unwrap().plugin = plugin.Merge(plugin.Collect(subT)...)

				subT.unwrap().plugin.Hooks.BeforeEach()
				defer subT.unwrap().plugin.Hooks.AfterEach()

				suiteHooks.BeforeEach(suiteClone, tt)
				defer suiteHooks.AfterEach(suiteClone, tt)

				test.Run(suite, subT)
			})
		}
	})
}

func Run[T commonT](t T, name string, f func(t T)) bool {
	name = t.unwrap().plugin.Plan.Rename(name)

	return t.Run(name, func(tt *testing.T) {
		subT := construct(&concreteT{T: tt}, &t)

		subT.unwrap().plugin = plugin.Merge(plugin.Collect(subT)...)

		subT.unwrap().plugin.Hooks.BeforeEach()
		defer subT.unwrap().plugin.Hooks.AfterEach()

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
	if value.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("expected value kind to be a pointer, got %s", value.Type()))
	}

	if value.Type() != parent.Type() {
		panic(fmt.Sprintf("value (%s) and parent (%s) type mismatch", value.Type(), parent.Type()))
	}

	if value.Type() == reflect.TypeOf(t) {
		value.Set(reflect.ValueOf(t))
		return
	}

	const initMethodName = "Init"

	initFunc := value.MethodByName(initMethodName)
	isPromoted := reflectutil.IsPromotedMethod(value.Type(), initMethodName)

	if initFunc.IsValid() && !isPromoted {
		initFuncType := initFunc.Type()

		isValidOut := initFuncType.NumOut() == 0
		isValidIn := initFuncType.NumIn() == 2 && initFuncType.In(0) == parent.Type()

		if !isValidIn || !isValidOut {
			t.Fatalf(
				"wrong signature for %[1sT.Init, must be: func (%[1]s) Init(%s, ...%s)",
				value.Type(), parent.Type(), reflect.TypeFor[plugin.Option](),
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
		valueField := value.Field(i)

		if !valueField.CanSet() {
			continue
		}

		var parentField reflect.Value

		if parent.IsValid() {
			parentField = parent.Field(i)
		} else {
			parentField = reflect.New(valueField.Type()).Elem()
		}

		if valueField.Kind() == reflect.Pointer {
			initValue(t, valueField, parentField, inits, options...)
			continue
		}

		initValue(t, valueField.Addr(), parentField.Addr(), inits, options...)
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
