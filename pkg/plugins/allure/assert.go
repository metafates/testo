package allure

import (
	"fmt"

	"github.com/metafates/testo"
)

//go:generate sh -c "cd _codegen && go run . -pkg allure -path ../assert.gen.go"

type CommonT interface {
	testo.CommonT

	Interface
}

func Require[T CommonT](t T) Requirements[T] {
	return Requirements[T]{t: t}
}

func Assert[T CommonT](t T) Assertions[T] {
	return Assertions[T]{t: t}
}

type Assertions[T CommonT] struct {
	t    T
	mode ParameterMode
}

type Requirements[T CommonT] struct {
	t    T
	mode ParameterMode
}

func (r Requirements[T]) Masked() Requirements[T] {
	r.mode = ParameterModeMasked

	return r
}

func (a Assertions[T]) Masked() Assertions[T] {
	a.mode = ParameterModeMasked

	return a
}

func (r Requirements[T]) Hidden() Requirements[T] {
	r.mode = ParameterModeHidden

	return r
}

func (a Assertions[T]) Hidden() Assertions[T] {
	a.mode = ParameterModeHidden

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
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
}
