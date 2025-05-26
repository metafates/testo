package testman

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"testman/internal/iterutil"
)

type suiteHooks[Suite any, T any] struct {
	BeforeAll  func(Suite, T)
	BeforeEach func(Suite, T)
	AfterEach  func(Suite, T)
	AfterAll   func(Suite, T)
}

type suiteTest[Suite any, T any] struct {
	Name string
	Run  func(Suite, T)
}

func collectSuiteHooks[Suite any, T fataller](t T) suiteHooks[Suite, T] {
	return suiteHooks[Suite, T]{
		BeforeAll:  getSuiteHook[Suite](t, "BeforeAll"),
		BeforeEach: getSuiteHook[Suite](t, "BeforeEach"),
		AfterEach:  getSuiteHook[Suite](t, "AfterEach"),
		AfterAll:   getSuiteHook[Suite](t, "AfterAll"),
	}
}

func collectSuiteTests[Suite any, T commonT](
	t *testing.T,
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

			tests = append(tests, suiteTest[Suite, T]{
				Name: name,
				Run: func(s Suite, t T) {
					casesValues := make(map[string][]reflect.Value, len(requiredCases))

					for name, c := range requiredCases {
						casesValues[name] = c.Func(s)
					}

					for i, params := range iterutil.Permutations(casesValues) {
						paramValue := reflect.New(param).Elem()

						paramsInterface := make(map[string]any, len(params))

						for name, value := range params {
							paramValue.FieldByName(name).Set(value)

							paramsInterface[name] = value.Interface()
						}

						// TODO: better name, compute %03d (e.g. %06d) from permutations count
						runSubtest(
							t,
							fmt.Sprintf("Case %03d", i),
							func(t T) { // init T
								t.unwrap().caseParams = paramsInterface
							},
							func(t T) { // actual test
								method.Func.Call([]reflect.Value{
									reflect.ValueOf(s),
									reflect.ValueOf(t),
									paramValue,
								})
							},
						)
					}
				},
			})

		default:
			wrongSignatureError()
		}
	}

	return tests
}

type suiteCase[Suite any] struct {
	Provides reflect.Type
	Func     func(Suite) []reflect.Value
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
