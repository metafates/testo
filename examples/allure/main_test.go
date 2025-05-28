package main_test

import (
	"testing"
	"time"

	"github.com/metafates/tego"
	"github.com/metafates/tego/pkg/plugins/allure"
	"github.com/metafates/tego/pkg/plugins/assertions"
)

type T struct {
	*tego.T

	// Плагины.
	// У плагинов есть свои хуки и возможность изменять стандартные методы типа Log, Error
	*assertions.Assertions
	*allure.Allure
}

func Test(t *testing.T) {
	tego.RunSuite[MySuite, *T](t, allure.WithOutputPath("allure-results"))
}

type MySuite struct{}

// Конечно, хуки типа BeforeAll, AfterEach доступны для имплементации
func (MySuite) BeforeEach(t *T) {
	tego.Run(t, "Connect to cosmos", func(t *T) {
		t.Log("Boop-beep...")
		time.Sleep(time.Second)

		tego.Run(t, "nested", func(t *T) { t.Log("works") })

		t.Log("Ready to test!")
	}, allure.AsSetup())

	tego.Run(t, "Another setup step", func(t *T) {
		time.Sleep(time.Second)
	}, allure.AsSetup())
}

func (MySuite) AfterEach(t *T) {
	tego.Run(t, "Say goodbye", func(t *T) {
		t.Log("Goodbye!")
	}, allure.AsTearDown())
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
	tego.Run(t, "ensure that value is true", func(t *T) {
		value := true

		// Эта функция исходит из assertions плагина
		t.Require().True(value)

		time.Sleep(time.Second)

		panic("oops")
	})

	tego.Run(t, "skip this step", func(t *T) {
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
