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

	"github.com/metafates/tego/constraint"
	"github.com/metafates/tego/internal/reflectutil"
	"github.com/metafates/tego/internal/stack"
	"github.com/metafates/tego/internal/suite"
	"github.com/metafates/tego/plugin"
)

const parallelWrapperTest = "!"

// RunSuite will run the tests under the given suite.
//
// It also accepts options for the plugins which can be used to configure those plugins.
// See [plugin.Option].
//
//nolint:thelper // not a helper
func RunSuite[Suite any, T constraint.T](t *testing.T, options ...plugin.Option) {
	suiteName := reflectutil.NameOf[Suite]()

	t.Run(suiteName, func(rawT *testing.T) {
		t := construct[T](rawT, nil, options...)
		unwrap(t, func(t *actualT) { t.suiteName = suiteName })

		runSuite[Suite](t)
	})
}

func runSuite[Suite any, T constraint.T](t T) {
	suiteHooks := suite.HooksOf[Suite](t)

	theSuite := reflectutil.Make[Suite]()

	unwrap(t, func(t *actualT) { t.plugin.Hooks.BeforeAll.Run() })
	suiteHooks.BeforeAll(theSuite, t)

	defer func() {
		suiteHooks.AfterAll(theSuite, t)
		unwrap(t, func(t *actualT) { t.plugin.Hooks.AfterAll.Run() })
	}()

	t.Run(parallelWrapperTest, func(rawT *testing.T) {
		runSuiteTests(construct(rawT, &t), theSuite, suiteHooks)
	})
}

func runSuiteTests[Suite any, T constraint.T](t T, s Suite, hooks suite.Hooks[Suite, T]) {
	tests := testsFor[Suite](t)

	for _, test := range tests {
		t.Run(test.Name, func(rawT *testing.T) {
			runSuiteTest(
				construct(rawT, &t),
				suite.Clone(s),
				hooks,
				test,
			)
		})
	}
}

func runSuiteTest[Suite any, T constraint.T](
	t T,
	s Suite,
	hooks suite.Hooks[Suite, T],
	test suite.Test[Suite, T],
) {
	unwrap(t, func(t *actualT) { t.plugin.Hooks.BeforeEach.Run() })
	hooks.BeforeEach(s, t)

	defer func() {
		hooks.AfterEach(s, t)
		unwrap(t, func(t *actualT) { t.plugin.Hooks.AfterEach.Run() })
	}()

	defer func() {
		if r := recover(); r != nil {
			unwrap(t, func(t *actualT) {
				t.panicInfo = &PanicInfo{
					Value: r,
					Trace: string(debug.Stack()),
				}
			})

			t.Errorf("Test %q panicked: %r", t.Name(), r)
		}
	}()

	test.Run(s, t)
}

func Run[T constraint.T](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	return runSubtest(t, name, nil, f, options...)
}

func runSubtest[T constraint.T](
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

		unwrap(subT, func(t *actualT) { t.plugin.Hooks.BeforeEach.Run() })

		// TODO: fix panic when running subtests inside cleanup.
		subT.Cleanup(func() {
			unwrap(subT, func(t *actualT) { t.plugin.Hooks.AfterEach.Run() })
		})

		defer func() {
			if r := recover(); r != nil {
				unwrap(subT, func(t *actualT) {
					t.panicInfo = &PanicInfo{
						Value: r,
						Trace: string(debug.Stack()),
					}
				})

				subT.Errorf("Test %q panicked: %v", subT.Name(), r)
			}
		}()

		subtest(subT)
	})
}

// construct will construct a new user T (inherits actual T)
// with the given parent and options.
func construct[T constraint.T](t *testing.T, parent *T, options ...plugin.Option) T {
	t.Helper()

	switch reflect.TypeFor[T]() {
	case reflect.TypeOf(t): // special case: T is actually *testing.T
		//nolint:forcetypeassert // checked with reflection
		return any(t).(T)

	case reflect.TypeFor[*actualT](): // special case: T is *tego.T
		actual := actualT{T: t, plugin: plugin.Merge(), levelOptions: options}

		if parent != nil {
			unwrap(*parent, func(t *actualT) { actual.parent = t })
		}

		//nolint:forcetypeassert // checked with reflection
		return any(&actual).(T)
	}

	value := reflectutil.Filled[T]()

	inits := stack.New[func()]()

	seedT := actualT{T: t, levelOptions: options}
	if parent != nil {
		unwrap(*parent, func(t *actualT) { seedT.parent = t })
	}

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
	// unwrap(value, func(t *actualT) {
	// 	t.plugin = plugin.Merge(plugin.Collect(&value)...)
	// 	t.levelOptions = options
	//
	// 	if parent != nil {
	// 		unwrap(*parent, func(parentT *actualT) {
	// 			t.parent = parentT
	// 		})
	// 	}
	// })

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

func testsFor[Suite any, T constraint.T](t T) []suite.Test[Suite, T] {
	cases := suite.CasesOf[Suite](t)
	tests := rawTestsFor(t, cases)

	unwrap(t, func(t *actualT) { tests = applyPlan(t.plugin.Plan, tests) })

	return tests
}

func rawTestsFor[Suite any, T constraint.T](
	t T,
	cases map[string]suite.Case[Suite],
) []suite.Test[Suite, T] {
	vt := reflect.TypeFor[Suite]()

	tests := make([]suite.Test[Suite, T], 0, vt.NumMethod())

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
			tests = append(tests, suite.Test[Suite, T]{
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

			tests = append(tests, newParametrizedTest[Suite, T](name, method, requiredCases))
		}
	}

	return tests
}

func newParametrizedTest[Suite any, T constraint.T](
	name string,
	method reflect.Method,
	cases map[string]suite.Case[Suite],
) suite.Test[Suite, T] {
	param := method.Type.In(2)

	return suite.Test[Suite, T]{
		Name: name,
		Run: func(s Suite, t T) {
			casesValues := make(map[string][]reflect.Value, len(cases))

			for name, c := range cases {
				casesValues[name] = c.Func(s)
			}

			var i int

			for params := range casesPermutations(casesValues) {
				i++

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
						unwrap(t, func(t *actualT) { t.caseParams = caseParams })
					},
					func(t T) { // actual test
						method.Func.Call([]reflect.Value{
							reflect.ValueOf(suite.Clone(s)),
							reflect.ValueOf(t),
							paramValue,
						})
					},
				)
			}
		},
	}
}

func applyPlan[Suite any, T constraint.T](
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
