package testo

import (
	"reflect"
	"strings"

	"github.com/metafates/testo/internal/pragma"
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
		Info TestInfo
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

func (plannedSuiteTest[Suite, T]) TestoInternal(pragma.DoNotImplement) {}

func (t plannedSuiteTest[Suite, T]) Name() string {
	return t.suiteTest.Name
}

func suiteCasesOf[Suite any, T fataller](t T) map[string]suiteCase[Suite] {
	vt := reflect.TypeFor[Suite]()

	cases := make(map[string]suiteCase[Suite])

	for i := range reflectutil.AsPointer(vt).NumMethod() {
		method := reflectutil.AsPointer(vt).Method(i)

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

	method, ok := reflectutil.AsPointer(suite).MethodByName(name)
	if !ok {
		return func(Suite, T) {}
	}

	switch f := method.Func.Interface().(type) {
	case func(Suite, T):
		return f

	case func(*Suite, T):
		return func(s Suite, t T) { f(&s, t) }

	default:
		t.Fatalf(
			"wrong signature for %[1]s.%[2]s, must be: func %[2]s(%T)",
			suite, name, t,
		)

		return nil
	}
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
