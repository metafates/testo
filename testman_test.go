package testman_test

import (
	"testing"

	"testman/plugin/allure"
	"testman/plugin/foo"

	"testman"
)

type T struct {
	*testman.T

	*foo.Foo
	*allure.Allure
}

func Test(t *testing.T) {
	testman.Run[Suite, T](t)
}

type Suite struct{}

func (s Suite) TestFoo(t *T) {
	t.Description("this is a sample test")

	testman.Subtest(t, "subtest", func(t *T) {
		t.Log("hi")
	})

	testman.Subtest(t, "weird", func(t *T) {
		t.Log("hello??")
	})
}
