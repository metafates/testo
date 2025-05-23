package testman

import (
	"testing"
)

type beforeAller[T testing.TB] interface {
	BeforeAll(t *T)
}

type beforeEacher[T testing.TB] interface {
	BeforeEach(t *T)
}

type afterEacher[T testing.TB] interface {
	AfterEach(t *T)
}

type afterAller[T testing.TB] interface {
	AfterAll(t *T)
}
