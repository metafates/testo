package reflectutil

import (
	"reflect"
	"testing"
)

func TestMakeValue(t *testing.T) {
	type Mock struct {
		String *string
		Nested *struct {
			Nested *struct {
				N *int
			}
		}
	}

	var value Mock

	FillValue(reflect.ValueOf(&value))

	notNil(t, value.Nested)
	notNil(t, value.String)
	notNil(t, value.Nested.Nested)
	notNil(t, value.Nested.Nested.N)
}

func notNil[T any](t *testing.T, value *T) {
	t.Helper()

	if value == nil {
		t.Fatal("nil value")
	}
}
