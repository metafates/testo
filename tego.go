// Package tego is a modular testing framework built on top of [testing.T].
package tego

import (
	"fmt"
	"iter"
	"maps"
	"reflect"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/metafates/tego/internal/reflectutil"
	"github.com/metafates/tego/internal/stack"
	"github.com/metafates/tego/internal/suite"
	"github.com/metafates/tego/plugin"
)

const parallelWrapperTest = "tego!"

// RunSuite will run the tests under the given suite.
//
// It also accepts options for the plugins which can be used to configure those plugins.
// See [plugin.Option].
//
//nolint:thelper // not a helper
func RunSuite[Suite any, T CommonT](t *testing.T, options ...plugin.Option) {
	suiteName := reflectutil.NameOf[Suite]()

	t.Run(suiteName, func(rawT *testing.T) {
		t := construct[T](rawT, nil, options...)
		t.unwrap().suiteName = suiteName

		runSuite[Suite](t)
	})
}

func runSuite[Suite any, T CommonT](t T) {
	suiteHooks := suite.HooksOf[Suite](t)

	theSuite := reflectutil.Make[Suite]()

	cases := suite.CasesOf[Suite](t)
	tests := testsFor(t, cases)

	t.unwrap().plugin.Hooks.BeforeAll.Run()
	suiteHooks.BeforeAll(theSuite, t)

	defer func() {
		suiteHooks.AfterAll(theSuite, t)
		t.unwrap().plugin.Hooks.AfterAll.Run()
	}()

	t.Run(parallelWrapperTest, func(rawT *testing.T) {
		tests := tests.All(suite.Clone(theSuite))
		tests = applyPlan(t.unwrap().plugin.Plan, tests)

		for _, test := range tests {
			rawT.Run(test.Name, func(rawT *testing.T) {
				runSuiteTest(
					construct(rawT, &t),
					suite.Clone(theSuite),
					suiteHooks,
					test,
				)
			})
		}
	})
}

func runSuiteTest[Suite any, T CommonT](
	t T,
	s Suite,
	hooks suite.Hooks[Suite, T],
	test suite.Test[Suite, T],
) {
	t.unwrap().caseParams = test.Params

	t.unwrap().plugin.Hooks.BeforeEach.Run()
	hooks.BeforeEach(s, t)

	defer func() {
		hooks.AfterEach(s, t)
		t.unwrap().plugin.Hooks.AfterEach.Run()
	}()

	defer func() {
		if r := recover(); r != nil {
			t.unwrap().panicInfo = &PanicInfo{
				Value: r,
				Trace: string(debug.Stack()),
			}

			t.Errorf("Test %q panicked: %r", t.Name(), r)
		}
	}()

	test.Run(s, t)
}

func Run[T CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	return runSubtest(t, name, nil, f, options...)
}

func runSubtest[T CommonT](
	tt T,
	name string,
	initT, subtest func(t T),
	options ...plugin.Option,
) bool {
	//nolint:thelper // not a helper
	return tt.Run(name, func(t *testing.T) {
		subT := construct(t, &tt, options...)

		if initT != nil {
			initT(subT)
		}

		subT.unwrap().plugin.Hooks.BeforeEach.Run()
		defer subT.unwrap().plugin.Hooks.AfterEach.Run()

		defer func() {
			if r := recover(); r != nil {
				subT.unwrap().panicInfo = &PanicInfo{
					Value: r,
					Trace: string(debug.Stack()),
				}

				subT.Errorf("Test %q panicked: %v", subT.Name(), r)
			}
		}()

		subtest(subT)
	})
}

// construct will construct a new user T (inherits actual T)
// with the given parent and options.
func construct[T CommonT](t *testing.T, parent *T, options ...plugin.Option) T {
	t.Helper()

	seedT := actualT{
		T:            t,
		levelOptions: options,
		plugin:       plugin.Merge(),
	}

	if parent != nil {
		seedT.parent = (*parent).unwrap()
	}

	switch reflect.TypeFor[T]() {
	case reflect.TypeFor[*actualT](): // special case: T is *tego.T
		//nolint:forcetypeassert // checked with reflection
		return any(&seedT).(T)
	}

	value := reflectutil.Filled[T]()

	inits := stack.New[func()]()

	initValue(
		&seedT,
		reflect.ValueOf(&value),
		reflect.ValueOf(parent),
		&inits,
	)

	for {
		init, ok := inits.Pop()
		if !ok {
			break
		}

		init()
	}

	seedT.plugin = plugin.Merge(plugin.Collect(&value)...)

	return value
}

func initValue(
	t *T,
	value, parent reflect.Value,
	inits *stack.Stack[func()],
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
				reflect.ValueOf(t.options()),
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
			initValue(t, valueField, parentField, inits)

			continue
		}

		initValue(t, valueField.Addr(), parentField.Addr(), inits)
	}
}

type suiteTests[Suite any, T CommonT] struct {
	Regular      []suite.Test[Suite, T]
	Parametrized []func(s Suite) []suite.Test[Suite, T]
}

func (st suiteTests[Suite, T]) All(s Suite) []suite.Test[Suite, T] {
	tests := st.Regular

	for _, p := range st.Parametrized {
		tests = append(tests, p(s)...)
	}

	return tests
}

func testsFor[Suite any, T CommonT](
	t T,
	cases map[string]suite.Case[Suite],
) suiteTests[Suite, T] {
	vt := reflect.TypeFor[Suite]()

	var tests suiteTests[Suite, T]

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
		default:
			wrongSignatureError()

		case 2: // regular test - (Suite, T)
			//nolint:forcetypeassert // checked by reflection
			tests.Regular = append(tests.Regular, suite.Test[Suite, T]{
				Name: name,
				Run:  method.Func.Interface().(func(Suite, T)),
			})

		case 3: // parametrized test - (Suite, T, Params)
			param := method.Type.In(2)

			requiredCases := make(map[string]suite.Case[Suite])

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

			tests.Parametrized = append(tests.Parametrized, newParametrizedTest[Suite, T](name, method, requiredCases))
		}
	}

	return tests
}

func newParametrizedTest[Suite any, T CommonT](
	name string,
	method reflect.Method,
	cases map[string]suite.Case[Suite],
) func(Suite) []suite.Test[Suite, T] {
	param := method.Type.In(2)

	return func(s Suite) []suite.Test[Suite, T] {
		casesValues := make(map[string][]reflect.Value, len(cases))

		for name, c := range cases {
			casesValues[name] = c.Func(s)
		}

		var i int

		var tests []suite.Test[Suite, T]

		for params := range casesPermutations(casesValues) {
			i++

			paramValue := reflect.New(param).Elem()

			caseParams := make(map[string]any, len(params))

			for name, value := range params {
				paramValue.FieldByName(name).Set(value)

				caseParams[name] = value.Interface()
			}
			tests = append(tests, suite.Test[Suite, T]{
				Name:   fmt.Sprintf("%s Case %d", name, i),
				Params: caseParams,
				Run: func(s Suite, t T) {
					method.Func.Call([]reflect.Value{
						reflect.ValueOf(suite.Clone(s)),
						reflect.ValueOf(t),
						paramValue,
					})
				},
			})
		}

		return tests
	}
}

func applyPlan[Suite any, T CommonT](
	plan plugin.Plan,
	tests []suite.Test[Suite, T],
) []suite.Test[Suite, T] {
	plannedTests := make([]plugin.PlannedTest, 0, len(tests))

	for _, t := range tests {
		plannedTests = append(plannedTests, t)
	}

	plannedTests = plan.Modify(plannedTests)

	testsToReturn := make([]suite.Test[Suite, T], 0, len(plannedTests))

	for _, t := range plannedTests {
		if t == nil {
			continue
		}

		st, ok := t.(suite.Test[Suite, T])
		if !ok {
			// TODO: better error message
			panic(fmt.Sprintf("test %q was modified", t.GetName()))
		}

		testsToReturn = append(testsToReturn, st)
	}

	return testsToReturn
}

// casesPermutations returns a determenistic permutations of the given cases values for test.
func casesPermutations[V any](v map[string][]V) iter.Seq[map[string]V] {
	keys := slices.Collect(maps.Keys(v))
	slices.Sort(keys)

	return func(yield func(map[string]V) bool) {
		current := make(map[string]V, len(keys))

		var walk func(i int) bool

		walk = func(i int) bool {
			if i == len(keys) {
				return yield(maps.Clone(current))
			}

			key := keys[i]
			for _, val := range v[key] {
				current[key] = val

				if !walk(i + 1) {
					return false
				}
			}

			return true
		}

		_ = walk(0)
	}
}
