package testman_test

import (
	"testing"
	"time"

	"testman/plugin/allure"
	"testman/plugin/foo"

	"testman"
)

type T struct {
	*testman.T

	foo.Foo
	*allure.Allure
}

func Test(t *testing.T) {
	testman.Run[Suite, T](t)
}

type Suite struct {
	number int
}

func (s *Suite) BeforeAll(t *T) {
	t.Log("Suite.BeforeAll")

	s.number = 5
}

func (s Suite) AfterAll(t *T) {
	t.Log("Suite.AfterAll")
}

func (s *Suite) BeforeEach(t *T) {
	t.Log("Suite.BeforeEach: " + t.Name())

	s.number *= 2
}

func (Suite) AfterEach(t *T) {
	t.Log("Suite.AfterEach: " + t.Name())
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

	t.Title("Example test")

	time.Sleep(2 * time.Second)
	t.Log(s.number)

	testman.Subtest(t, "subtest here!", func(t *T) {
		defer t.Measure()()

		t.Log("Hello!!!")
		time.Sleep(time.Second)
	})
}
