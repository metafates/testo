package testo

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/metafates/testo/internal/directive"
	"github.com/metafates/testo/internal/maputil"
	"github.com/metafates/testo/internal/reflectutil"
	"github.com/metafates/testo/plugin"
)

type (
	suiteHooks[Suite any, T any] struct {
		BeforeAll  func(Suite, T)
		BeforeEach func(Suite, T)
		AfterEach  func(Suite, T)
		AfterAll   func(Suite, T)
	}

	suiteTest[Suite any, T any] struct {
		Name string
		Info plugin.TestInfo
		Run  func(Suite, T)
	}

	suiteCase[Suite any] struct {
		Provides reflect.Type
		Func     func(Suite) []reflect.Value
	}
)

var _ plugin.PlannedTest = (*plannedSuiteTest[any, any])(nil)

type plannedSuiteTest[Suite any, T any] struct {
	suiteTest[Suite, T]
}

func (plannedSuiteTest[Suite, T]) TestoInternal(directive.DoNotImplement) {}

func (t plannedSuiteTest[Suite, T]) Name() string {
	return t.suiteTest.Name
}

func (t plannedSuiteTest[Suite, T]) Info() plugin.TestInfo {
	return t.suiteTest.Info
}

// isTest states whether name is a valid test name (or other type, according to prefix).
//
// It checks if the next character after prefix is uppercase.
//
//	TestFoo    => true
//	Test       => true
//	TestfooBar => false
func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}

	// "Test" is ok
	if len(name) == len(prefix) {
		return true
	}

	r, _ := utf8.DecodeRuneInString(name[len(prefix):])

	return !unicode.IsLower(r)
}

func suiteCasesOf[Suite any, T fataller](t T) map[string]suiteCase[Suite] {
	vt := reflect.TypeFor[Suite]()

	cases := make(map[string]suiteCase[Suite])

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		const prefix = "Cases"

		if !isTest(method.Name, prefix) {
			continue
		}

		name := strings.TrimPrefix(method.Name, prefix)

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
				var suite reflect.Value

				if method.Type.In(0).Kind() == reflect.Pointer &&
					reflect.TypeOf(s).Kind() != reflect.Pointer {
					suite = reflect.ValueOf(&s)
				} else {
					suite = reflect.ValueOf(s)
				}

				slice := method.Func.Call([]reflect.Value{suite})[0]

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

// suiteHooksOf returns hooks of the given suite.
func suiteHooksOf[Suite any, T fataller](t T) suiteHooks[Suite, T] {
	return suiteHooks[Suite, T]{
		BeforeAll:  getHook[Suite](t, "BeforeAll"),
		BeforeEach: getHook[Suite](t, "BeforeEach"),
		AfterEach:  getHook[Suite](t, "AfterEach"),
		AfterAll:   getHook[Suite](t, "AfterAll"),
	}
}

type fataller interface {
	Fatalf(format string, args ...any)
}

func getHook[Suite any, T fataller](t T, name string) func(Suite, T) {
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

		return nil
	}

	return f
}

// cloner can clone itself.
type cloner[Self any] interface {
	// Clone returns a new instance cloned from the caller.
	Clone() Self
}

func cloneSuite[Suite any](suite Suite) Suite {
	if cloner, ok := any(suite).(cloner[Suite]); ok {
		return cloner.Clone()
	}

	return reflectutil.DeepClone(suite)
}

// suiteTests contains all the suite tests.
//
// While regular tests are ready to be run,
// parametrized tests are tricky.
//
// We can't know how many permutations (hence number of tests)
// they will have until we receive all values for each case by calling CasesXXX funcs.
// However, we can't do that before running the BeforeAll hooks - cases funcs may
// depend on in being run first.
//
// But we should not run any hooks until we are sure that tests are correct
// and no error should be raised.
//
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

//nolint:cyclop,funlen // splitting it would make it even more complex
func testsFor[Suite any, T CommonT](
	t T,
	cases map[string]suiteCase[Suite],
) suiteTests[Suite, T] {
	vt := reflect.TypeFor[Suite]()

	var tests suiteTests[Suite, T]

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}

		raiseWrongSignatureError := func() {
			t.Fatalf(
				"wrong signature for %[1]s.%[2]s, must be: func %[1]s.%[2]s(%[3]s) or func %[1]s.%[2]s(%[3]s, struct{...})",
				vt,
				method.Name,
				reflect.TypeFor[T](),
			)
		}

		if method.Type.NumOut() != 0 {
			raiseWrongSignatureError()
		}

		if method.Type.NumIn() < 2 {
			raiseWrongSignatureError()
		}

		if method.Type.In(1) != reflect.TypeFor[T]() {
			raiseWrongSignatureError()
		}

		if method.Type.NumIn() == 3 && method.Type.In(2).Kind() != reflect.Struct {
			raiseWrongSignatureError()
		}

		switch method.Type.NumIn() {
		default:
			raiseWrongSignatureError()

		case 2: // regular test - (Suite, T)
			//nolint:forcetypeassert // checked by reflection
			tests.Regular = append(tests.Regular, suiteTest[Suite, T]{
				Name: method.Name,
				Info: plugin.RegularTestInfo{
					RawBaseName: method.Name,
					Level:       1,
				},
				Run: method.Func.Interface().(func(Suite, T)),
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
				newParametrizedTest[Suite, T](method.Name, method, requiredCases),
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

		for _, params := range casesPermutations(casesValues) {
			i++

			paramValue := reflect.New(param).Elem()

			caseParams := make(map[string]any, len(params))

			for name, value := range params {
				paramValue.FieldByName(name).Set(value)

				caseParams[name] = value.Interface()
			}

			tests = append(tests, suiteTest[Suite, T]{
				// TODO: better name? Allow plugins customize it?
				Name: fmt.Sprintf("%s case %d", name, i),
				Info: plugin.ParametrizedTestInfo{
					RawBaseName: name,
					Params:      caseParams,
				},
				Run: func(s Suite, t T) {
					method.Func.Call([]reflect.Value{
						reflect.ValueOf(s),
						reflect.ValueOf(t),
						paramValue,
					})
				},
			})
		}

		return tests
	}
}

// casesPermutations returns a determenistic permutations of the given cases values for test.
func casesPermutations[V any](v map[string][]V) []map[string]V {
	var result []map[string]V

	keys := maputil.Keys(v)

	// Sort keys for consistent processing order (optional but ensures deterministic output)
	slices.Sort(keys)

	var generatePermutations func(current map[string]V, index int)

	generatePermutations = func(current map[string]V, index int) {
		// Base case: if all keys have been processed
		if index == len(keys) {
			result = append(result, maps.Clone(current))

			return
		}

		key := keys[index]

		for _, val := range v[key] {
			current[key] = val
			generatePermutations(current, index+1)
		}
	}

	current := make(map[string]V)

	generatePermutations(current, 0)

	return result
}
