package allure

import (
	"fmt"

	"github.com/metafates/testo"
)

//go:generate sh -c "cd _codegen && go run . -pkg allure -path ../assert.gen.go"

// CommonT is interface which
// all T's with Allure plugin installed implement.
type CommonT interface {
	testo.CommonT

	Interface
}

// Require returns a new [Requirements] instance.
func Require[T CommonT](t T) Requirements[T] {
	return Requirements[T]{t: t}
}

// Assert returns a new [Assertions] instance.
func Assert[T CommonT](t T) Assertions[T] {
	return Assertions[T]{t: t}
}

// Assertions provides a set of helpers to perform assertions in tests.
//
// Each assertion is included in the Allure report
// as a step with passed parameters.
type Assertions[T CommonT] struct {
	t    T
	mode ParameterMode
}

// Requirements implements the same assertions as [Assertions]
// but stops test execution when assertion fails.
type Requirements[T CommonT] struct {
	t    T
	mode ParameterMode
}

// Masked returns a new requirements instance
// which will mask its parameters.
func (r Requirements[T]) Masked() Requirements[T] {
	r.mode = ParameterModeMasked

	return r
}

// Masked returns a new assertions instance
// which will mask its parameters.
func (a Assertions[T]) Masked() Assertions[T] {
	a.mode = ParameterModeMasked

	return a
}

func messageFromMsgAndArgs(msgAndArgs ...any) string {
	switch len(msgAndArgs) {
	case 0:
		return ""

	case 1:
		msg := msgAndArgs[0]

		if s, ok := msg.(string); ok {
			return s
		}

		return fmt.Sprintf("%+v", msg)

	default:
		format, ok := msgAndArgs[0].(string)
		if !ok {
			panic(fmt.Sprintf("format must be a string, got %T", msgAndArgs[0]))
		}

		return fmt.Sprintf(format, msgAndArgs[1:]...)
	}
}
