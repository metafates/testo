package main

import (
	"testing"

	"testman"
)

type T struct {
	*testman.T

	PluginWhichReversesTestOrder
	PluginWhichAddsNewTests
	PluginWhichSkipsRandomTests
	PluginWhichOverridesLog
}

func Test(t *testing.T) {
	testman.Suite[Suite, *T](t)
}

type Suite struct{}

func (Suite) Test1(t *T) {
	t.Log("Hello from Test1")
}

func (Suite) Test2(t *T) {
	t.Log("Hello from Test2")
}

func (Suite) Test3(t *T) {
	t.Log("Hello from Test3")
}
