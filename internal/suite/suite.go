package suite

import (
	"reflect"
	"strings"

	"github.com/metafates/tego/internal/reflectutil"
	"github.com/metafates/tego/plugin"
)

type (
	Hooks[Suite any, T any] struct {
		BeforeAll  func(Suite, T)
		BeforeEach func(Suite, T)
		AfterEach  func(Suite, T)
		AfterAll   func(Suite, T)
	}

	Test[Suite any, T any] struct {
		Name string
		Info plugin.TestInfo
		Run  func(Suite, T)
	}

	Case[Suite any] struct {
		Provides reflect.Type
		Func     func(Suite) []reflect.Value
	}
)

func (t Test[Suite, T]) GetName() string {
	return t.Name
}

func CasesOf[Suite any, T fataller](t T) map[string]Case[Suite] {
	vt := reflect.TypeFor[Suite]()

	cases := make(map[string]Case[Suite])

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

		cases[name] = Case[Suite]{
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

// HooksOf returns hooks of the given suite.
func HooksOf[Suite any, T fataller](t T) Hooks[Suite, T] {
	return Hooks[Suite, T]{
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
	}

	return f
}

// cloner can clone itself.
type cloner[Self any] interface {
	// Clone returns a new instance cloned from the caller.
	Clone() Self
}

func Clone[Suite any](suite Suite) Suite {
	if cloner, ok := any(suite).(cloner[Suite]); ok {
		return cloner.Clone()
	}

	return reflectutil.DeepClone(suite)
}
