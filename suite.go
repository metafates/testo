package testman

import (
	"reflect"
	"strings"
	"testing"
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

func collectSuiteTests[Suite any, T fataller](t *testing.T) []suiteTest[Suite, T] {
	vt := reflect.TypeFor[Suite]()

	tests := make([]suiteTest[Suite, T], 0, vt.NumMethod())

	for i := range vt.NumMethod() {
		method := vt.Method(i)

		if !method.IsExported() {
			continue
		}

		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}

		switch f := method.Func.Interface().(type) {
		case func(Suite, T):
			tests = append(tests, suiteTest[Suite, T]{
				Name: method.Name,
				Run:  f,
			})

		default:
			t.Fatalf(
				"wrong signature for %[1]s.%[2]s, must be: func %[1]s.%[2]s(t %s)",
				reflect.TypeFor[Suite](),
				method.Name,
				reflect.TypeFor[T](),
			)
		}
	}

	return tests
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
