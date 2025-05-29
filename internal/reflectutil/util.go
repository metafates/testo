package reflectutil

import (
	"fmt"
	"reflect"
	"strings"
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

// NameOf returns name of the underlying type T.
func NameOf[T any]() string {
	t := reflect.TypeFor[T]()

	return Elem(t).Name()
}

// Make a new zero value of T.
//
// As a special case for pointers it will
// return pointer to the zero value of T (not nil).
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

// DeepClone returns a deep clone of v.
func DeepClone[T any](v T) T {
	cloneMap := make(map[uintptr]reflect.Value)

	cloned, ok := deepClone(reflect.ValueOf(v), cloneMap).Interface().(T)
	if !ok {
		panic(fmt.Sprintf("cloned value must be of the same type as original (%T)", v))
	}

	return cloned
}

func deepClone(v reflect.Value, cloneMap map[uintptr]reflect.Value) reflect.Value {
	if !v.IsValid() {
		return reflect.Value{}
	}

	switch v.Kind() {
	case reflect.String:
		//nolint:forcetypeassert // checked with kind
		return reflect.ValueOf(strings.Clone(v.Interface().(string)))

	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return v

	case reflect.Array:
		cloned := reflect.New(reflect.ArrayOf(v.Len(), v.Type().Elem())).Elem()

		for i := range v.Len() {
			cloned.Index(i).Set(deepClone(v.Index(i), cloneMap))
		}

		return cloned

	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		cloned := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := range v.Len() {
			cloned.Index(i).Set(deepClone(v.Index(i), cloneMap))
		}

		return cloned

	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		cloned := reflect.MakeMap(v.Type())

		for _, key := range v.MapKeys() {
			clonedKey := deepClone(key, cloneMap)
			clonedValue := deepClone(v.MapIndex(key), cloneMap)
			cloned.SetMapIndex(clonedKey, clonedValue)
		}

		return cloned

	case reflect.Struct:
		cloned := reflect.New(v.Type()).Elem()

		for i := range v.NumField() {
			if v.Type().Field(i).PkgPath == "" { // Exported field
				cloned.Field(i).Set(deepClone(v.Field(i), cloneMap))
			}
		}

		return cloned

	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		ptr := v.Pointer()
		if cloned, ok := cloneMap[ptr]; ok {
			return cloned
		}

		clonedElem := deepClone(v.Elem(), cloneMap)
		clonedPtr := reflect.New(v.Type().Elem())
		clonedPtr.Elem().Set(clonedElem)
		cloneMap[ptr] = clonedPtr

		return clonedPtr

	case reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		return deepClone(v.Elem(), cloneMap)

	default:
		panic("unsupported kind: " + v.Kind().String())
	}
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
