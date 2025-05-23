package testman_test

import (
	"fmt"
	"testing"

	"testman/plugin/foo"

	"testman"
)

type T struct {
	*testman.T

	foo.Foo
}

func (t *T) Run(name string, f func(t *T)) bool {
	return t.T.Run(name, func(t *testman.T) {
		f(&T{T: t})
	})
}

func Test(t *testing.T) {
	tt := T{
		T: testman.New(t),
	}

	testman.Run[Suite](&tt)
}

type Suite struct {
	number int
}

func (s *Suite) BeforeAll(t *T) {
	t.Log("before all once")
	s.number = 5
}

func (s Suite) AfterAll(t *T) {
	t.Log("Done")
}

func (s *Suite) BeforeEach(t *T) {
	s.number *= 2
	fmt.Println("Running before " + t.Name())
}

func (Suite) AfterEach(t *T) {
	fmt.Println("Running after " + t.Name())
}

func (s Suite) TestBar(t *T) {
	t.Parallel()

	t.Log(s.number)
}

func (s Suite) TestFoo(t *T) {
	t.Label("name")
	t.Parallel()

	t.Log(s.number)

	t.Run("nested", func(t *T) {
		t.Log("hello from " + t.Name())

		t.Run("even more", func(t *T) {
			t.Log("hello again but from " + t.Name())
		})
	})
}
