package main

import (
	"testing"

	"github.com/metafates/tego"
)

type T = tego.T

func Test(t *testing.T) {
	tego.RunSuite[Suite, *T](t)
}

type Suite struct{}

func (Suite) AfterAll(t *T) {
	t.Log(">>>>>> suite tear down", t.Name())
}

func (Suite) AfterEach(t *T) {
	t.Log(">>>>>> single test tear down", t.Name(), t.Panicked())
}

func (Suite) TestOne(t *T) {
	t.Parallel()

	for _, v := range []string{"sub1", "sub2", "sub3"} {
		tego.Run(t, v, func(t *T) {
			t.Parallel()
			t.Log("panicking now")
			panic("oh no")
		})
	}
}

func (Suite) TestTwo(t *T) {
	t.Parallel()

	for _, v := range []string{"sub1", "sub2", "sub3"} {
		tego.Run(t, v, func(t *T) { t.Parallel() })
	}
}

func (Suite) TestThree(t *T) {
	t.Parallel()

	panic("ooooops")
}
