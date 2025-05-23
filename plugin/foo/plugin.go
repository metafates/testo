package foo

import "testman"

type Foo struct{}

func (Foo) New(t *testman.T) Foo {
	return Foo{}
}
