package main

import (
	"testing"
	"time"

	"github.com/metafates/tego"
)

type T struct {
	*tego.T

	PluginWhichReversesTestOrder
	PluginWhichAddsNewTests
	PluginWhichSkipsRandomTests
	PluginWhichOverridesLog
	PluginWhichRenamesTests

	Timer
}

func Test(t *testing.T) {
	tego.RunSuite[Suite, *T](t)
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

	time.Sleep(time.Second)
}
