package testman_test

import (
	"fmt"
	"testing"
	"time"

	"testman/plugin/foo"

	"testman"
)

type T struct {
	*testman.T

	foo.Foo
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

	defer t.Measure()()

	t.Log(s.number)

	testman.Subtest(t, "assertions", func(t *T) {
		t.RequireTrue(false)
	})
}

func (s Suite) TestFoo(t *T) {
	t.Parallel()

	defer t.Measure()()

	time.Sleep(2 * time.Second)
	t.Log(s.number)

	testman.Subtest(t, "subtest here!", func(t *T) {
		defer t.Measure()()

		fmt.Println("Hello!!!")
		time.Sleep(time.Second)
	})
}
