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

	value := Filled[Mock]()

	notNil(t, value.Nested)
	notNil(t, value.String)
	notNil(t, value.Nested.Nested)
	notNil(t, value.Nested.Nested.N)
}

func TestDeepClone(t *testing.T) {
	type Mock struct {
		Name    string
		private []string
	}

	original := Mock{Name: "Test", private: []string{"test"}}

	clone := DeepClone(original)

	equal(t, original, clone)
}

func notNil[T any](t *testing.T, value *T) {
	t.Helper()

	if value == nil {
		t.Fatal("nil value")
	}
}

func equal(t *testing.T, want, got any) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("not equal: want %v, got %v", want, got)
	}
}
