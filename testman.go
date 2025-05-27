package testman

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"slices"
	"testing"

	"testman/internal/reflectutil"
	"testman/internal/stack"
	"testman/plugin"
)

// RunSuite will run the tests under the given suite.
//
// It also accepts options for the plugins which can be used to configure those plugins.
// See [plugin.Option].
func RunSuite[Suite any, T commonT](t *testing.T, options ...plugin.Option) {
	tt := construct[T](&concreteT{T: t}, nil, options...)
	tt.unwrap().suiteName = reflect.TypeFor[Suite]().Name()

	cases := collectSuiteCases[Suite](tt)
	tests := applyPlan(tt.unwrap().plugin.Plan, collectSuiteTests[Suite, T](t, cases))

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	suiteHooks := collectSuiteHooks[Suite](tt)

	var suite Suite

	tt.unwrap().plugin.Hooks.BeforeAll.Run()
	suiteHooks.BeforeAll(suite, tt)

	defer func() {
		suiteHooks.AfterAll(suite, tt)
		tt.unwrap().plugin.Hooks.AfterAll.Run()
	}()

	// wrap all tests so that AfterAll hooks will
	// be called after these tests even if they use Parallel().
	t.Run(tt.unwrap().SuiteName(), func(t *testing.T) {
		for _, test := range tests {
			var suiteClone Suite

			if cloner, ok := any(suite).(Cloner[Suite]); ok {
				suiteClone = cloner.Clone()
			} else {
				suiteClone = reflectutil.DeepClone(suite)
			}

			t.Run(test.Name, func(t *testing.T) {
				subT := construct(&concreteT{T: t}, &tt)

				subT.unwrap().plugin.Hooks.BeforeEach.Run()
				suiteHooks.BeforeEach(suiteClone, tt)

				defer func() {
					if r := recover(); r != nil {
						subT.unwrap().panicInfo = &PanicInfo{
							Msg:   r,
							Trace: string(debug.Stack()),
						}

						subT.Fail()
					}

					suiteHooks.AfterEach(suiteClone, tt)
					subT.unwrap().plugin.Hooks.AfterEach.Run()
				}()

				test.Run(suite, subT)
			})
		}
	})
}

func Run[T commonT](t T, name string, f func(t T)) bool {
	return runSubtest(t, name, nil, f)
}

func runSubtest[T commonT](tt T, name string, initT, subtest func(t T)) bool {
	name = tt.unwrap().plugin.Plan.Rename(name)

	return tt.Run(name, func(t *testing.T) {
		subT := construct(&concreteT{T: t}, &tt)

		if initT != nil {
			initT(subT)
		}

		subT.unwrap().plugin.Hooks.BeforeEach.Run()

		defer func() {
			if r := recover(); r != nil {
				subT.unwrap().panicInfo = &PanicInfo{
					Msg:   r,
					Trace: string(debug.Stack()),
				}

				subT.Fail()
			}

			subT.unwrap().plugin.Hooks.AfterEach.Run()
		}()

		subtest(subT)
	})
}

func construct[T commonT](t *concreteT, parent *T, options ...plugin.Option) T {
	value := reflectutil.Filled[T]()

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

	value.unwrap().plugin = plugin.Merge(plugin.Collect(&value)...)

	if parent != nil {
		value.unwrap().parent = (*parent).unwrap()
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
