package main_test

import (
	"testing"
	"time"

	"testman/plugin/allure"
	"testman/plugin/assertions"

	"testman"
)

type T struct {
	*testman.T

	// Плагины.
	// У плагинов есть свои хуки и возможность изменять стандартные методы типа Log, Error
	*assertions.Assertions
	*allure.Allure
}

func Test(t *testing.T) {
	testman.Suite[MySuite, T](t, allure.WithOutputPath("my-allure-results"))
}

type MySuite struct{}

// Конечно, хуки типа BeforeAll, AfterEach доступны для имплементации
func (s MySuite) BeforeAll(t T) {
	t.Logf("BeforeAll")
}

func (s MySuite) AfterAll(t T) {
	t.Logf("AfterAll ")
}

func (s MySuite) TestFoo(t T) {
	// Параллельные тесты поддерживаются
	t.Parallel()

	// Эти функции исходят из allure плагина
	t.Description("this is a sample test")
	t.Labels(allure.Label{Name: "tag", Value: "Q924"})
	t.Links(allure.Link{Name: "github", URL: "https://github.com", Type: "issue"})
	t.Flaky()

	// Плагин Allure превратит эти подтесты в шаги в репорте
	testman.Run(t, "ensure that value is true", func(t T) {
		value := true

		// Эта функция исходит из assertions плагина
		t.Require().True(value)

		time.Sleep(time.Second)
	})

	testman.Run(t, "skip this step", func(t T) {
		time.Sleep(2 * time.Second)

		t.Skip("skipped")
	})
}

func (s MySuite) TestAnotherParallel(t T) {
	t.Parallel()

	t.Require().True(true)
	time.Sleep(time.Second)
}
