package reflectutil

import (
	"reflect"
)

type canElem[Self any] interface {
	Kind() reflect.Kind
	Elem() Self
}

// Elem unwraps the underlying elem of the pointer.
//
// Nested pointers are also supported - e.g. given "****value" it will return "value".
//
// Non-pointer values will be returned as is.
func Elem[T canElem[T]](v T) T {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return v
}

// IsPromotedMethod reports whether method "name" on type "t" comes
// from an embedded (anonymous) field.
func IsPromotedMethod(t reflect.Type, name string) bool {
	// make sure we detect methods on either T or *T
	if _, ok := t.MethodByName(name); !ok {
		_, ok = reflect.PointerTo(t).MethodByName(name)
		if !ok {
			return false
		}
	}

	// now scan anonymous fields for a method of that name
	return walkEmbedded(t, name, make(map[reflect.Type]struct{}))
}

func Name[T any]() string {
	t := reflect.TypeFor[T]()

	return Elem(t).Name()
}

func Make[T any]() T {
	t := reflect.TypeFor[T]()

	var zero T

	if t.Kind() == reflect.Pointer {
		elem := reflect.ValueOf(&zero).Elem()

		elem.Set(reflect.New(t.Elem()))
	}

	return zero
}

// Filled returns a new value T with all the exported pointer fields recursively set to non-nil zero values.
// That is, if T is a struct and contains field *int it will be set to &0.
// That logic is also applies for nested structs.
func Filled[T any]() T {
	var value T

	fillValue(reflect.ValueOf(&value))

	return value
}

func fillValue(v reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		fillValue(v.Elem())

	case reflect.Struct:
		for i := range v.NumField() {
			field := v.Field(i)

			if field.CanSet() {
				fillValue(field)
			}
		}
	}
}

func walkEmbedded(t reflect.Type, name string, seen map[reflect.Type]struct{}) bool {
	t = Elem(t)

	if _, ok := seen[t]; ok {
		return false
	}

	seen[t] = struct{}{}

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := range t.NumField() {
		f := t.Field(i)
		if !f.Anonymous {
			continue
		}

		ft := Elem(f.Type)

		// check value‐receiver methods on ft
		if _, ok := ft.MethodByName(name); ok {
			return true
		}

		// check pointer‐receiver methods on *ft
		if _, ok := reflect.PointerTo(ft).MethodByName(name); ok {
			return true
		}

		// recurse through deeper anonymous embeddings
		if walkEmbedded(ft, name, seen) {
			return true
		}
	}

	return false
}
