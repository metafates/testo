package order_test

import (
	"testing"

	"testman"
)

type T struct {
	*testman.T
}

func Test(t *testing.T) {
	testman.Suite[Suite, T](t)
}

type Suite struct{}

func (Suite) TestFirst(t *T) {
	t.Log("I am the first!")
}

func (Suite) TestSecond(t *T) {
	t.Log("I am the second!")
}
