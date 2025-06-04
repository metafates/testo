package main

import (
	"testing"

	"github.com/metafates/testo"
)

type T = testo.T

func Test(t *testing.T) {
	testo.RunSuite[Suite, *T](t)
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
		testo.Run(t, v, func(t *T) {
			t.Parallel()
		})
	}
}

func (Suite) TestTwo(t *T) {
	t.Parallel()

	for _, v := range []string{"sub1", "sub2", "sub3"} {
		testo.Run(t, v, func(t *T) { t.Parallel() })
	}
}
