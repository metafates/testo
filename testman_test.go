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

	*foo.Foo
	*allure.Allure
}

func Test(t *testing.T) {
	testman.Run[Suite, T](t)
}

type Suite struct{}

func (s Suite) TestFoo(t *T) {
	t.Description("this is a sample test")
	t.Labels(allure.Label{Name: "tag", Value: "Q924"})
	t.Links(allure.Link{Name: "github", URL: "https://github.com", Type: "tms"})
	t.Flaky()

	testman.Subtest(t, "subtest", func(t *T) {
		t.Log("hi")

		time.Sleep(2 * time.Second)
	})

	testman.Subtest(t, "skipped", func(t *T) {
		t.Log("hello??")

		t.Skip()
	})

	testman.Subtest(t, "another one", func(t *T) {
		t.Log("ok")

		testman.Subtest(t, "nested", func(t *T) {
			t.Log("nested!")
		})
	})
}
