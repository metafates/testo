// Package testo is a modular testing framework built on top of [testing.T].
package testo

import (
	"fmt"
	"reflect"
	"runtime/debug"
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
// Suite type must be a pointer in a form of *MySuite.
//
// It also accepts options for the plugins which can be used to configure those plugins.
// See [plugin.Option].
func RunSuite[Suite any, T CommonT](t *testing.T, options ...plugin.Option) {
	t.Helper()

	if !reflectutil.IsSinglePointer(reflect.TypeFor[Suite]()) {
		panic(fmt.Sprintf(
			"invalid suite type specified '%s', did you mean '*%[2]s'?",
			reflect.TypeFor[Suite](),
			reflectutil.Elem(reflect.TypeFor[Suite]()),
		))
	}

	suiteName := reflectutil.NameOf[Suite]()

	options = append(getDefaultOptions(), options...)

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
	t.Helper()

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

	t.unwrap().T.Run(parallelWrapperTest, func(rawT *testing.T) {
		tests := tests.Get(cloneSuite(suite))
		tests = applyPlan(t.unwrap().plugin.Plan, tests)

		for _, test := range tests {
			rawT.Run(test.Name, func(rawT *testing.T) {
				innerT := construct(
					rawT,
					&t,
					func(t *actualT) {
						t.info.Test = test.Info
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
			t.unwrap().info.Panic = &plugin.PanicInfo{
				Value: r,
				Trace: string(debug.Stack()),
			}

			t.Fatalf("Test %q panicked: %r", t.Name(), r)
		}
	}()

	test.Run(s, t)
}

// Run a subtest.
// It has the same purpose as [testing.T.Run] but
// retains the passed [T] type for the subtest function.
func Run[T CommonT](
	t T,
	name string,
	f func(t T),
	options ...plugin.Option,
) bool {
	t.Helper()

	parentT := t

	return parentT.unwrap().T.Run(name, func(tt *testing.T) {
		t := construct(
			tt,
			&parentT,
			func(t *actualT) {
				t.info.Test = plugin.RegularTestInfo{
					RawBaseName: name,
					Level:       t.level(),
				}
			},
			options...,
		)

		t.unwrap().plugin.Hooks.BeforeEachSub.Run()
		defer t.unwrap().plugin.Hooks.AfterEachSub.Run()

		defer func() {
			if r := recover(); r != nil {
				t.unwrap().info.Panic = &plugin.PanicInfo{
					Value: r,
					Trace: string(debug.Stack()),
				}

				t.Errorf("test %q panicked: %v", t.Name(), r)
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

	seedT.info.Plugins = plugins
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

//nolint:cyclop,funlen // splitting it would make it even more complex
func initValue(
	t *T,
	value, parent reflect.Value,
	inits *stack.Stack[func()],
) {
	t.Helper()

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

func applyPlan[Suite any, T CommonT](
	plan plugin.Plan,
	tests []suiteTest[Suite, T],
) []suiteTest[Suite, T] {
	plannedTests := make([]plugin.PlannedTest, 0, len(tests))

	for _, t := range tests {
		plannedTests = append(plannedTests, plannedSuiteTest[Suite, T]{t})
	}

	plan.Modify(&plannedTests)

	testsToReturn := make([]suiteTest[Suite, T], 0, len(plannedTests))

	for _, t := range plannedTests {
		if t == nil {
			continue
		}

		planned, ok := t.(plannedSuiteTest[Suite, T])
		if !ok {
			// TODO: better error message, but it must be unreachable.
			panic(fmt.Sprintf("type of test %q is not suiteTest", t.Name()))
		}

		testsToReturn = append(testsToReturn, planned.suiteTest)
	}

	return testsToReturn
}
