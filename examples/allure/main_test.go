package main_test

import (
	"testing"
	"time"

	"testman"
	"testman/pkg/plugins/allure"
	"testman/pkg/plugins/assertions"
)

type T struct {
	*testman.T

	// Плагины.
	// У плагинов есть свои хуки и возможность изменять стандартные методы типа Log, Error
	*assertions.Assertions
	*allure.Allure
}

func Test(t *testing.T) {
	testman.RunSuite[MySuite, *T](t, allure.WithOutputPath("allure-results"))
}

type MySuite struct{}

// Конечно, хуки типа BeforeAll, AfterEach доступны для имплементации
func (s MySuite) BeforeAll(t *T) {
	t.Logf("BeforeAll")
}

func (s MySuite) AfterAll(t *T) {
	t.Logf("AfterAll")
}

func (MySuite) TestFoo(t *T) {
	// Параллельные тесты поддерживаются
	t.Parallel()

	// Эти функции исходят из allure плагина
	t.Description("This is a sample test")
	t.Labels(allure.Label{Name: "tag", Value: "Q924"})
	t.Links(allure.Link{Name: "github", URL: "https://github.com", Type: "issue"})
	t.Flaky()

	// Плагин Allure превратит эти подтесты в шаги в репорте
	testman.Run(t, "ensure that value is true", func(t *T) {
		value := true

		// Эта функция исходит из assertions плагина
		t.Require().True(value)

		time.Sleep(time.Second)

		panic("oops")
	})

	testman.Run(t, "skip this step", func(t *T) {
		time.Sleep(2 * time.Second)

		t.Skip("skipped")
	})
}

func (MySuite) TestAnotherParallel(t *T) {
	t.Parallel()

	t.Require().True(true)
	time.Sleep(time.Second)
}

func (MySuite) CasesName() []string {
	return []string{"John", "Jane", "Bob"}
}

func (MySuite) CasesAge() []int {
	return []int{5, 20, 40, 70}
}

func (MySuite) TestBar(t *T, params struct {
	Name string
	Age  int
},
) {
	t.Require().True(len(params.Name) > 0)
	t.Require().True(params.Age > 18)
}
