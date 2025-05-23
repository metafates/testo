package testman_test

import (
	"fmt"
	"testing"

	"testman/plugin/foo"

	"testman"
)

// TODO: initialize T using reflection (call New() for each field that implements it)
// If not - use default values. User can implement its own New function, but its optional

type T struct {
	*testman.T

	foo.Foo
}

func (T) New(tt *testman.T) T {
	return T{T: tt}
}

func (t *T) Run(name string, f func(t *T)) bool {
	return t.T.Run(name, func(t *testman.T) {
		f(&T{T: t})
	})
}

func Test(t *testing.T) {
	testman.Run[Suite, T](t)
}

type Suite struct {
	number int
}

func (s *Suite) BeforeAll(t *T) {
	t.Log("BeforeAll")

	s.number = 5
}

func (s Suite) AfterAll(t *T) {
	t.Log("AfterAll")
}

func (s *Suite) BeforeEach(t *T) {
	fmt.Println("BeforeEach: " + t.Name())

	s.number *= 2
}

func (Suite) AfterEach(t *T) {
	fmt.Println("AfterEach: " + t.Name())
}

func (s Suite) TestBar(t *T) {
	t.Parallel()

	t.Log(s.number)
}

func (s Suite) TestFoo(t *T) {
	t.Parallel()

	t.Log(s.number)

	t.Run("nested", func(t *T) {
		t.Log("hello from " + t.Name())

		t.Run("even more", func(t *T) {
			t.Log("hello again but from " + t.Name())
		})
	})
}
