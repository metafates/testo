package reflectutil

import "reflect"

// Elem unwraps the underlying elem of the pointer.
//
// Nested pointers are also supported - e.g. given "****value" it will return "value".
//
// Non-pointer values will be returned as is.
func Elem(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return v
}
