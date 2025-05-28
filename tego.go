package tego

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/metafates/tego/internal/iterutil"
	"github.com/metafates/tego/internal/reflectutil"
	"github.com/metafates/tego/internal/stack"
	"github.com/metafates/tego/plugin"
)

// RunSuite will run the tests under the given suite.
//
// It also accepts options for the plugins which can be used to configure those plugins.
// See [plugin.Option].
//
//nolint:thelper // not a helper
func RunSuite[Suite any, T commonT](t *testing.T, options ...plugin.Option) {
	tt := construct[T](&concreteT{T: t}, nil, options...)
	tt.unwrap().suiteName = reflect.TypeFor[Suite]().Name()

	tests := testsFor[Suite](tt)

	// nothing to do
	if len(tests) == 0 {
		t.Log("warn: no tests to run")

		return
	}

	suiteHooks := collectHooksFor[Suite](tt)

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
			suiteClone := cloneSuite(suite)

			t.Run(test.Name, func(t *testing.T) {
				subT := construct(&concreteT{T: t}, &tt)

				subT.unwrap().plugin.Hooks.BeforeEach.Run()
				suiteHooks.BeforeEach(suiteClone, subT)

				defer func() {
					if r := recover(); r != nil {
						subT.unwrap().panicInfo = &PanicInfo{
							Msg:   r,
							Trace: string(debug.Stack()),
						}

						subT.Fail()
					}

					suiteHooks.AfterEach(suiteClone, subT)
					subT.unwrap().plugin.Hooks.AfterEach.Run()
				}()

				test.Run(suite, subT)
			})
		}
	})
}

func Run[T commonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	return runSubtest(t, name, nil, f, options...)
}

func runSubtest[T commonT](
	tt T,
	name string,
	initT, subtest func(t T),
	options ...plugin.Option,
) bool {
	name = tt.unwrap().plugin.Plan.Rename(name)

	//nolint:thelper // not a helper
	return tt.Run(name, func(t *testing.T) {
		subT := construct(&concreteT{T: t}, &tt, options...)

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

type (
	suiteHooks[Suite any, T any] struct {
		BeforeAll  func(Suite, T)
		BeforeEach func(Suite, T)
		AfterEach  func(Suite, T)
		AfterAll   func(Suite, T)
	}

	suiteTest[Suite any, T any] struct {
		Name string
		Run  func(Suite, T)
	}

	suiteCase[Suite any] struct {
		Provides reflect.Type
		Func     func(Suite) []reflect.Value
	}
)

// construct will construct a new user T (inherits actual T)
// with the given parent and options.
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
				"wrong signature for %[1]s.Init, must be: func (%[1]s) Init(%s, ...%s)",
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

func collectHooksFor[Suite any, T fataller](t T) suiteHooks[Suite, T] {
	return suiteHooks[Suite, T]{
		BeforeAll:  getSuiteHook[Suite](t, "BeforeAll"),
		BeforeEach: getSuiteHook[Suite](t, "BeforeEach"),
		AfterEach:  getSuiteHook[Suite](t, "AfterEach"),
		AfterAll:   getSuiteHook[Suite](t, "AfterAll"),
	}
}

func rawTestsFor[Suite any, T commonT](
	t T,
	cases map[string]suiteCase[Suite],
) []suiteTest[Suite, T] {
	vt := reflect.TypeFor[Suite]()

	tests := make([]suiteTest[Suite, T], 0, vt.NumMethod())

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		name, ok := strings.CutPrefix(method.Name, "Test")
		if !ok {
			continue
		}

		wrongSignatureError := func() {
			t.Fatalf(
				"wrong signature for %[1]s.%[2]s, must be: func %[1]s.%[2]s(%[3]s) or func %[1]s.%[2]s(%[3]s, struct{...})",
				reflect.TypeFor[Suite](),
				method.Name,
				reflect.TypeFor[T](),
			)
		}

		if method.Type.NumOut() != 0 {
			wrongSignatureError()
		}

		if method.Type.NumIn() < 2 {
			wrongSignatureError()
		}

		if method.Type.In(0) != vt || method.Type.In(1) != reflect.TypeFor[T]() {
			wrongSignatureError()
		}

		if method.Type.NumIn() == 3 && method.Type.In(2).Kind() != reflect.Struct {
			wrongSignatureError()
		}

		switch method.Type.NumIn() {
		case 2: // regular test - (Suite, T)
			//nolint:forcetypeassert // checked by reflection
			tests = append(tests, suiteTest[Suite, T]{
				Name: name,
				Run:  method.Func.Interface().(func(Suite, T)),
			})

		case 3: // parametrized test - (Suite, T, Params)
			param := method.Type.In(2)

			requiredCases := make(map[string]suiteCase[Suite])

			for i := range param.NumField() {
				field := param.Field(i)

				c, ok := cases[field.Name]
				if !ok {
					t.Fatalf(
						"wrong param signature for %[1]s.%[2]s: Cases%[3]s for param %[3]q not found",
						reflect.TypeFor[Suite](),
						method.Name,
						field.Name,
					)
				}

				if !c.Provides.AssignableTo(field.Type) {
					// TODO: "of type ..." shows invalid type
					t.Fatalf(
						"wrong param signature for %[1]s.%[2]s: Cases%[3]s provides %s values, not assignable to param %[3]q of type %s",
						reflect.TypeFor[Suite](),
						method.Name,
						field.Name,
						c.Provides,
						field.Type,
					)
				}

				requiredCases[field.Name] = c
			}

			tests = append(tests, parametrizedSuiteTest[Suite, T](name, method, requiredCases))

		default:
			wrongSignatureError()
		}
	}

	return tests
}

func parametrizedSuiteTest[Suite any, T commonT](
	name string,
	method reflect.Method,
	cases map[string]suiteCase[Suite],
) suiteTest[Suite, T] {
	param := method.Type.In(2)

	return suiteTest[Suite, T]{
		Name: name,
		Run: func(s Suite, t T) {
			casesValues := make(map[string][]reflect.Value, len(cases))

			for name, c := range cases {
				casesValues[name] = c.Func(s)
			}

			for i, params := range iterutil.Permutations(casesValues) {
				paramValue := reflect.New(param).Elem()

				caseParams := make(map[string]any, len(params))

				for name, value := range params {
					paramValue.FieldByName(name).Set(value)

					caseParams[name] = value.Interface()
				}

				runSubtest(
					t,
					fmt.Sprintf("Case %d", i),
					func(t T) { // init T
						t.unwrap().caseParams = caseParams
					},
					func(t T) { // actual test
						method.Func.Call([]reflect.Value{
							reflect.ValueOf(cloneSuite(s)),
							reflect.ValueOf(t),
							paramValue,
						})
					},
				)
			}
		},
	}
}

func collectSuiteCases[Suite any, T fataller](t T) map[string]suiteCase[Suite] {
	vt := reflect.TypeFor[Suite]()

	cases := make(map[string]suiteCase[Suite])

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		name, ok := strings.CutPrefix(method.Name, "Cases")
		if !ok {
			continue
		}

		isValidIn := method.Type.NumIn() == 1
		isValidOut := method.Type.NumOut() == 1 && method.Type.Out(0).Kind() == reflect.Slice

		if !isValidIn || !isValidOut {
			t.Fatalf(
				"wrong signature for %[1]s.%[2]s, must be: func (%[1]s) %[2]s() []...",
				reflect.TypeFor[Suite](), method.Name, t,
			)
		}

		cases[name] = suiteCase[Suite]{
			Provides: method.Type.Out(0).Elem(),
			Func: func(s Suite) []reflect.Value {
				slice := method.Func.Call([]reflect.Value{reflect.ValueOf(s)})[0]

				values := make([]reflect.Value, 0, slice.Len())

				for i := range slice.Len() {
					v := slice.Index(i)

					values = append(values, v)
				}

				return values
			},
		}
	}

	return cases
}

func getSuiteHook[Suite any, T fataller](t T, name string) func(Suite, T) {
	suite := reflect.TypeFor[Suite]()

	method, ok := suite.MethodByName(name)
	if !ok {
		return func(Suite, T) {}
	}

	f, ok := method.Func.Interface().(func(Suite, T))
	if !ok {
		t.Fatalf(
			"wrong signature for %[1]s.%[2]s, must be: func %[2]s(%T)",
			suite, name, t,
		)
	}

	return f
}

func cloneSuite[T any](suite T) T {
	if cloner, ok := any(suite).(Cloner[T]); ok {
		return cloner.Clone()
	}

	return reflectutil.DeepClone(suite)
}

func testsFor[Suite any, T commonT](t T) []suiteTest[Suite, T] {
	cases := collectSuiteCases[Suite](t)
	tests := rawTestsFor(t, cases)

	return applyPlan(t.unwrap().plugin.Plan, tests)
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
