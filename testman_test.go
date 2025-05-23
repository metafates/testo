package testman_test

import (
	"fmt"
	"testing"

	"testman/plugin/foo"
	"testman/plugin/myplugin"

	"testman"
)

type Plugins interface {
	myplugin.Plugin
	foo.Plugin
}

func Test(t *testing.T) {
	testman.Run[Plugins](t, new(Suite), myplugin.New())
}

type Suite struct{}

func (Suite) Test(t testman.T) {
	fmt.Println("hi!")

	t.Run("nested", func(t testman.T) {
		fmt.Println("hi 2!")
	})
}

func (Suite) TestSomething(t testman.TP[Plugins]) {
	t.P().Label("name")

	t.RunP("nested with plugin", func(t testman.TP[Plugins]) {
		t.P().Label("wow nested?")
	})

	fmt.Println("Hello!")
}
