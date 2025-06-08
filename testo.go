// Package testo is a modular testing framework built on top of [testing.T].
package testo

import (
	"fmt"
	"iter"
	"maps"
	"reflect"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/metafates/testo/internal/reflectutil"
	"github.com/metafates/testo/internal/stack"
	"github.com/metafates/testo/plugin"
)

// parallelWrapperTest is the name of tests which
// wrap multiple (possibly parallel) tests to ensure
// hooks are executed properly.
//
// It should contain some special symbol like exclamation mark,
// so that it won't collide with suite type name.
const parallelWrapperTest = "testo!"

// RunSuite will run the tests under the given suite.
//
// It also accepts options for the plugins which can be used to configure those plugins.
// See [plugin.Option].
//
//nolint:thelper // not a helper
func RunSuite[Suite any, T CommonT](t *testing.T, options ...plugin.Option) {
	suiteName := reflectutil.NameOf[Suite]()

	t.Run(suiteName, func(rawT *testing.T) {
		t := construct[T](
			rawT,
			nil,
			func(t *actualT) {
				t.suiteName = suiteName
			},
			options...,
		)

		runSuite[Suite](t)
	})
}

func runSuite[Suite any, T CommonT](t T) {
	suiteHooks := suiteHooksOf[Suite](t)

	suite := reflectutil.Make[Suite]()

	cases := suiteCasesOf[Suite](t)
	tests := testsFor(t, cases)

	t.unwrap().plugin.Hooks.BeforeAll.Run()
	suiteHooks.BeforeAll(suite, t)

	defer func() {
		suiteHooks.AfterAll(suite, t)
		t.unwrap().plugin.Hooks.AfterAll.Run()
	}()

	//nolint:thelper // naming this rawT makes this more readable.
	t.unwrap().T.Run(parallelWrapperTest, func(rawT *testing.T) {
		tests := tests.Get(cloneSuite(suite))
		tests = applyPlan(t.unwrap().plugin.Plan, tests)

		for _, test := range tests {
			rawT.Run(test.Name, func(rawT *testing.T) {
				innerT := construct(
					rawT,
					&t,
					func(t *actualT) {
						t.extra.Test = test.Info
					},
				)

				runSuiteTest(
					innerT,
					cloneSuite(suite),
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
	hooks suiteHooks[Suite, T],
	test suiteTest[Suite, T],
) {
	t.unwrap().plugin.Hooks.BeforeEach.Run()
	hooks.BeforeEach(s, t)

	defer func() {
		hooks.AfterEach(s, t)
		t.unwrap().plugin.Hooks.AfterEach.Run()
	}()

	defer func() {
		if r := recover(); r != nil {
			t.unwrap().extra.Panic = &PanicInfo{
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
	parentT := t

	//nolint:thelper // not a helper
	return parentT.unwrap().T.Run(name, func(tt *testing.T) {
		t := construct(
			tt,
			&parentT,
			func(t *actualT) {
				t.extra.Test = RegularTestInfo{RawBaseName: name}
			},
			options...,
		)

		t.unwrap().plugin.Hooks.BeforeEach.Run()
		defer t.unwrap().plugin.Hooks.AfterEach.Run()

		defer func() {
			if r := recover(); r != nil {
				t.unwrap().extra.Panic = &PanicInfo{
					Value: r,
					Trace: string(debug.Stack()),
				}

				t.Errorf("Test %q panicked: %v", t.Name(), r)
			}
		}()

		f(t)
	})
}

// construct will construct a new user T (inherits actual T)
// with the given parent and options.
func construct[T CommonT](
	t *testing.T,
	parent *T,
	fill func(t *actualT),
	options ...plugin.Option,
) T {
	t.Helper()

	seedT := actualT{
		T:            t,
		levelOptions: options,
		plugin:       plugin.MergeSpecs(),
	}

	if parent != nil {
		seedT.parent = (*parent).unwrap()
	}

	if fill != nil {
		fill(&seedT)
	}

	// special case: T is *testo.T
	if reflect.TypeFor[T]() == reflect.TypeFor[*actualT]() {
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

	// inits are deferred because we should run Init only
	// when all the fields are ready.
	for {
		init, ok := inits.Pop()
		if !ok {
			break
		}

		init()
	}

	plugins := plugin.Collect(&value)

	seedT.extra.Plugins = plugins
	seedT.plugin = mergePlugins(plugins...)

	return value
}

func mergePlugins(plugins ...plugin.Plugin) plugin.Spec {
	specs := make([]plugin.Spec, 0, len(plugins))

	for _, p := range plugins {
		specs = append(specs, p.Plugin())
	}

	return plugin.MergeSpecs(specs...)
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
		if !value.CanAddr() {
			// TODO: add path to the field so that it is clear where error happens
			panic(fmt.Sprintf("using non-pointer value of %s", reflect.TypeFor[T]()))
		}

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

// suiteTests contains all the suite tests.
//
// While regular tests are ready to be run,
// parametrized tests are tricky.
// We can't know how many permutations (hence number of tests)
// they will have until we all values for each case by calling CasesXXX funcs.
// However, we can't do that before running the BeforeAll hooks,
// since it would confuse users and make it less useful overall.
// But we should not run any hooks until we are sure that tests are correct
// and no error should be raised (static analysis).
// That's why we statically analyze parametrized tests signatures,
// but delay the actual collection for later.
type suiteTests[Suite any, T CommonT] struct {
	Regular      []suiteTest[Suite, T]
	Parametrized []func(s Suite) []suiteTest[Suite, T]
}

// Get all suite tests.
//
// Suite instance is required here to get
// parameter cases (CasesXXX funcs), not to invoke the actual tests.
func (st suiteTests[Suite, T]) Get(s Suite) []suiteTest[Suite, T] {
	tests := st.Regular

	for _, p := range st.Parametrized {
		tests = append(tests, p(s)...)
	}

	return tests
}

func testsFor[Suite any, T CommonT](
	t T,
	cases map[string]suiteCase[Suite],
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
			tests.Regular = append(tests.Regular, suiteTest[Suite, T]{
				Name: name,
				Info: RegularTestInfo{RawBaseName: name},
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

			tests.Parametrized = append(
				tests.Parametrized,
				newParametrizedTest[Suite, T](name, method, requiredCases),
			)
		}
	}

	return tests
}

func newParametrizedTest[Suite any, T CommonT](
	name string,
	method reflect.Method,
	cases map[string]suiteCase[Suite],
) func(Suite) []suiteTest[Suite, T] {
	param := method.Type.In(2)

	return func(s Suite) []suiteTest[Suite, T] {
		casesValues := make(map[string][]reflect.Value, len(cases))

		for name, c := range cases {
			casesValues[name] = c.Func(s)
		}

		var (
			tests []suiteTest[Suite, T]
			i     int
		)

		for params := range casesPermutations(casesValues) {
			i++

			paramValue := reflect.New(param).Elem()

			caseParams := make(map[string]any, len(params))

			for name, value := range params {
				paramValue.FieldByName(name).Set(value)

				caseParams[name] = value.Interface()
			}

			tests = append(tests, suiteTest[Suite, T]{
				Name: fmt.Sprintf("%s case %d", name, i),
				Info: ParametrizedTestInfo{
					BaseName: name,
					Params:   caseParams,
				},
				Run: func(s Suite, t T) {
					method.Func.Call([]reflect.Value{
						reflect.ValueOf(cloneSuite(s)),
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
	tests []suiteTest[Suite, T],
) []suiteTest[Suite, T] {
	plannedTests := make([]plugin.PlannedTest, 0, len(tests))

	for _, t := range tests {
		plannedTests = append(plannedTests, t)
	}

	plannedTests = plan.Modify(plannedTests)

	testsToReturn := make([]suiteTest[Suite, T], 0, len(plannedTests))

	for _, t := range plannedTests {
		if t == nil {
			continue
		}

		st, ok := t.(suiteTest[Suite, T])
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
